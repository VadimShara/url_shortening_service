package service_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/VadimShara/url_shortening_service/internal/service"
	"github.com/VadimShara/url_shortening_service/internal/tests/mocks"
	"github.com/VadimShara/url_shortening_service/pkg/errs"
)

type Test struct {
	name         string
	verifyResult func(t *testing.T, s *service.Url, testName string)
}

func Test_Url_SaveUrl(t *testing.T) {
	testData := map[string][]string{
		"valid": {
			"https://example.com",
			"https://google.com",
		},
		"invalid": {
			"invalid-url",
			"",
		},
	}

	expectedData := map[string][]string{
		"valid": {
			"alias1",
			"alias2",
		},
		"invalid": {},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUrlSaver := mocks.NewMockUrlSaver(ctrl)
	mockUrlRedirecter := mocks.NewMockUrlRedirecter(ctrl)

	s := service.New(slog.Default(), mockUrlSaver, mockUrlRedirecter)
	ctx := context.Background()

	mockUrlSaver.EXPECT().SaveUrl(ctx, "https://example.com", gomock.Any()).Return("alias1", nil).AnyTimes()
	mockUrlSaver.EXPECT().SaveUrl(ctx, "https://google.com", gomock.Any()).Return("alias2", nil).AnyTimes()
	mockUrlSaver.EXPECT().SaveUrl(ctx, "invalid-url", gomock.Any()).Return("", errors.New("invalid url")).AnyTimes()
	mockUrlSaver.EXPECT().SaveUrl(ctx, "", gomock.Any()).Return("", errors.New("empty url")).AnyTimes()

	tests := []Test{
		{
			name: "valid",
			verifyResult: func(t *testing.T, s *service.Url, testName string) {
				for i := range testData[testName] {
					alias, err := s.SaveUrl(ctx, testData[testName][i])
					assert.NoError(t, err)
					assert.Equal(t, expectedData[testName][i], alias)
				}
			},
		},
		{
			name: "invalid",
			verifyResult: func(t *testing.T, s *service.Url, testName string) {
				for i := range testData[testName] {
					alias, err := s.SaveUrl(ctx, testData[testName][i])
					assert.Error(t, err)
					assert.Empty(t, alias)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.verifyResult(t, s, test.name)
		})
	}
}

func Test_Url_RedirectUrl(t *testing.T) {
	testData := map[string][]string{
		"valid": {
			"alias1",
			"alias2",
		},
		"invalid": {
			"unknown-alias",
			"",
		},
	}

	expectedData := map[string][]string{
		"valid": {
			"https://example.com",
			"https://google.com",
		},
		"invalid": {},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUrlSaver := mocks.NewMockUrlSaver(ctrl)
	mockUrlRedirecter := mocks.NewMockUrlRedirecter(ctrl)

	s := service.New(slog.Default(), mockUrlSaver, mockUrlRedirecter)
	ctx := context.Background()

	mockUrlRedirecter.EXPECT().GetUrl(ctx, "alias1").Return("https://example.com", nil).AnyTimes()
	mockUrlRedirecter.EXPECT().GetUrl(ctx, "alias2").Return("https://google.com", nil).AnyTimes()
	mockUrlRedirecter.EXPECT().GetUrl(ctx, "unknown-alias").Return("", errs.ErrUrlNotFound).AnyTimes()
	mockUrlRedirecter.EXPECT().GetUrl(ctx, "").Return("", errors.New("empty alias")).AnyTimes()

	tests := []Test{
		{
			name: "valid",
			verifyResult: func(t *testing.T, s *service.Url, testName string) {
				for i := range testData[testName] {
					url, err := s.RedirectUrl(ctx, testData[testName][i])
					assert.NoError(t, err)
					assert.Equal(t, expectedData[testName][i], url)
				}
			},
		},
		{
			name: "invalid",
			verifyResult: func(t *testing.T, s *service.Url, testName string) {
				for i := range testData[testName] {
					url, err := s.RedirectUrl(ctx, testData[testName][i])
					assert.Error(t, err)
					assert.Empty(t, url)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.verifyResult(t, s, test.name)
		})
	}
}
