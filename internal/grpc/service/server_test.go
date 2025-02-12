package urlgrpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	service "github.com/VadimShara/url_shortening_service/api/gen/go"
)

type mockUrl struct {
	SaveUrlFunc     func(ctx context.Context, url string) (string, error)
	RedirectUrlFunc func(ctx context.Context, alias string) (string, error)
}

func (m *mockUrl) SaveUrl(ctx context.Context, url string) (string, error) {
	return m.SaveUrlFunc(ctx, url)
}

func (m *mockUrl) RedirectUrl(ctx context.Context, alias string) (string, error) {
	return m.RedirectUrlFunc(ctx, alias)
}

func Test_SaveUrl(t *testing.T) {
	testData := map[string]string{
		"valid":   "https://example.com",
		"invalid": "",
	}

	expectedData := map[string]interface{}{
		"valid":   &service.SaveUrlResponse{Alias: "abcde12345"},
		"invalid": codes.InvalidArgument,
	}

	tests := []struct {
		name         string
		setupMock    func(*mockUrl)
		request      *service.SaveUrlRequest
		verifyResult func(t *testing.T, resp *service.SaveUrlResponse, err error, testName string)
	}{
		{
			name: "valid",
			setupMock: func(m *mockUrl) {
				m.SaveUrlFunc = func(ctx context.Context, url string) (string, error) {
					return "abcde12345", nil
				}
			},
			request: &service.SaveUrlRequest{Url: testData["valid"]},
			verifyResult: func(t *testing.T, resp *service.SaveUrlResponse, err error, testName string) {
				assert.NoError(t, err)
				assert.Equal(t, expectedData[testName], resp)
			},
		},
		{
			name:      "invalid",
			setupMock: func(m *mockUrl) {},
			request:   &service.SaveUrlRequest{Url: testData["invalid"]},
			verifyResult: func(t *testing.T, resp *service.SaveUrlResponse, err error, testName string) {
				assert.Error(t, err)
				assert.Equal(t, expectedData[testName], status.Code(err))
			},
		},
	}

	for _, test := range tests {
		mock := &mockUrl{}
		test.setupMock(mock)
		server := &serverApi{url: mock}

		t.Run(test.name, func(t *testing.T) {
			resp, err := server.SaveUrl(context.Background(), test.request)
			test.verifyResult(t, resp, err, test.name)
		})
	}
}

func Test_RedirectUrl(t *testing.T) {
	testData := map[string]string{
		"valid":   "abcde12345",
		"invalid": "",
	}

	expectedData := map[string]interface{}{
		"valid":   &service.RedirectUrlResponse{Url: "https://example.com"},
		"invalid": codes.InvalidArgument,
	}

	tests := []struct {
		name         string
		setupMock    func(*mockUrl)
		request      *service.RedirectUrlRequest
		verifyResult func(t *testing.T, resp *service.RedirectUrlResponse, err error, testName string)
	}{
		{
			name: "valid",
			setupMock: func(m *mockUrl) {
				m.RedirectUrlFunc = func(ctx context.Context, alias string) (string, error) {
					return "https://example.com", nil
				}
			},
			request: &service.RedirectUrlRequest{Alias: testData["valid"]},
			verifyResult: func(t *testing.T, resp *service.RedirectUrlResponse, err error, testName string) {
				assert.NoError(t, err)
				assert.Equal(t, expectedData[testName], resp)
			},
		},
		{
			name:      "invalid",
			setupMock: func(m *mockUrl) {},
			request:   &service.RedirectUrlRequest{Alias: testData["invalid"]},
			verifyResult: func(t *testing.T, resp *service.RedirectUrlResponse, err error, testName string) {
				assert.Error(t, err)
				assert.Equal(t, expectedData[testName], status.Code(err))
			},
		},
	}

	for _, test := range tests {
		mock := &mockUrl{}
		test.setupMock(mock)
		server := &serverApi{url: mock}

		t.Run(test.name, func(t *testing.T) {
			resp, err := server.RedirectUrl(context.Background(), test.request)
			test.verifyResult(t, resp, err, test.name)
		})
	}
}
