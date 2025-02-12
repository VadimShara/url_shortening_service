package postgres

import (
	"context"
	"fmt"
	"testing"

	errs "github.com/VadimShara/url_shortening_service/pkg/errs"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	//"github.com/jackc/pgx/v5/pgxpool"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPool struct {
	mock.Mock
}

func (m *MockPool) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	args = append([]interface{}{ctx, query}, args...)
	m.Called(args...)
	return nil // Мокируем возвращаемое значение для QueryRow
}

func (m *MockPool) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	args = append([]interface{}{ctx, query}, args...)
	m.Called(args...)
	return nil, nil // Мокируем возвращаемое значение для Exec
}

func TestSaveUrl(t *testing.T) {
	tests := []struct {
		name          string
		urlToSave     string
		alias         string
		mockResponse  string
		mockError     error
		expectedErr   error
		expectedAlias string
	}{
		{
			name:          "Save url successfully",
			urlToSave:     "https://test.com",
			alias:         "test",
			mockResponse:  "",
			mockError:     nil,
			expectedErr:   nil,
			expectedAlias: "test",
		},
		{
			name:          "Url already exists",
			urlToSave:     "https://test.com",
			alias:         "test",
			mockResponse:  "test",
			mockError:     nil,
			expectedErr:   errs.ErrUrlExists,
			expectedAlias: "test",
		},
		{
			name:          "Error during URL check",
			urlToSave:     "https://test.com",
			alias:         "test",
			mockResponse:  "",
			mockError:     fmt.Errorf("some database error"),
			expectedErr:   fmt.Errorf("repo.postgres.SaveUrl: failed to check for existing url: some database error"),
			expectedAlias: "",
		},
		{
			name:          "Error inserting URL",
			urlToSave:     "https://test.com",
			alias:         "test",
			mockResponse:  "",
			mockError:     fmt.Errorf("some insert error"),
			expectedErr:   fmt.Errorf("repo.postgres.SaveUrl: failed to save url and alias: some insert error"),
			expectedAlias: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := mocks.NewMockFakePoolInterface(ctrl)
			dummyError := fmt.Errorf("db.BeginTx failed with err: connection failed")
			pool.EXPECT().BeginTx(ctx, pgx.TxOptions{}).Return(nil, dummyError)

			mockPool.On("QueryRow", mock.Anything, "SELECT alias FROM url WHERE url = $1", tt.urlToSave).Return(tt.mockResponse, tt.mockError)
			mockPool.On("Exec", mock.Anything, "INSERT INTO url (url, alias) VALUES ($1, $2)", tt.urlToSave, tt.alias).Return(nil, tt.mockError)

			alias, err := db.SaveUrl(context.Background(), tt.urlToSave, tt.alias)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedAlias, alias)
		})
	}
}

func TestGetUrl(t *testing.T) {
	tests := []struct {
		name         string
		alias        string
		mockResponse string
		mockError    error
		expectedErr  error
		expectedUrl  string
	}{
		{
			name:         "Get url successfully",
			alias:        "test",
			mockResponse: "https://test.com",
			mockError:    nil,
			expectedErr:  nil,
			expectedUrl:  "https://test.com",
		},
		{
			name:         "Alias not found",
			alias:        "nonexistent",
			mockResponse: "",
			mockError:    fmt.Errorf("no rows"),
			expectedErr:  fmt.Errorf("repo.postgres.GetUrl: url not found"),
			expectedUrl:  "",
		},
		{
			name:         "Error during url retrieval",
			alias:        "test",
			mockResponse: "",
			mockError:    fmt.Errorf("some database error"),
			expectedErr:  fmt.Errorf("repo.postgres.GetUrl: some database error"),
			expectedUrl:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPool := new(MockPool)
			db := &DB{Pool: mockPool}

			mockPool.On("QueryRow", mock.Anything, "SELECT url FROM url WHERE alias = $1", tt.alias).Return(tt.mockResponse, tt.mockError)

			url, err := db.GetUrl(context.Background(), tt.alias)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedUrl, url)
		})
	}
}
