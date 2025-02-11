package inMemory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VadimShara/url_shortening_service/internal/repo/inMemory"
	"github.com/VadimShara/url_shortening_service/pkg/errs"
)

func Test_DB_SaveUrl(t *testing.T) {
	testData := map[string][]struct {
		url   string
		alias string
		err   error
	}{
		"valid": {
			{"https://example.com", "qwertyuiop", nil},
			{"https://google.com", "abcdefghij", nil},
		},
		"invalid": {
			{"https://example.com", "xxxxxxxxxx", errs.ErrUrlExists},
			{"https://newsite.com", "abcdefghij", errs.ErrAliasExists},
		},
	}

	tests := []struct {
		name         string
		verifyResult func(t *testing.T, db *inMemory.DB, testName string)
	}{
		{
			name: "valid",
			verifyResult: func(t *testing.T, db *inMemory.DB, testName string) {
				for _, data := range testData[testName] {
					alias, err := db.SaveUrl(context.Background(), data.url, data.alias)
					assert.ErrorIs(t, err, data.err)
					if err == nil {
						assert.Equal(t, data.alias, alias)
					}
				}
			},
		},
		{
			name: "invalid",
			verifyResult: func(t *testing.T, db *inMemory.DB, testName string) {
				db.SaveUrl(context.Background(), "https://example.com", "qwertyuiop")
				db.SaveUrl(context.Background(), "https://google.com", "abcdefghij")

				for _, data := range testData[testName] {
					_, err := db.SaveUrl(context.Background(), data.url, data.alias)
					assert.ErrorIs(t, err, data.err)
				}
			},
		},
	}

	for _, test := range tests {
		db := inMemory.NewDB()
		t.Run(test.name, func(t *testing.T) {
			test.verifyResult(t, db, test.name)
		})
	}
}

func TestGetUrl(t *testing.T) {
	testData := map[string][]struct {
		alias       string
		expectedUrl string
		expectedErr error
	}{
		"valid": {
			{"qwertyuiop", "https://example.com", nil},
			{"abcdefghij", "https://google.com", nil},
		},
		"invalid": {
			{"unknown", "", errs.ErrUrlNotFound},
		},
	}

	tests := []struct {
		name         string
		prepare      func(db *inMemory.DB)
		verifyResult func(t *testing.T, db *inMemory.DB, testName string)
	}{
		{
			name: "valid",
			prepare: func(db *inMemory.DB) {
				db.SaveUrl(context.Background(), "https://example.com", "qwertyuiop")
				db.SaveUrl(context.Background(), "https://google.com", "abcdefghij")
			},
			verifyResult: func(t *testing.T, db *inMemory.DB, testName string) {
				for _, data := range testData[testName] {
					url, err := db.GetUrl(context.Background(), data.alias)
					assert.NoError(t, err)
					assert.Equal(t, data.expectedUrl, url)
				}
			},
		},
		{
			name:    "invalid",
			prepare: func(db *inMemory.DB) {},
			verifyResult: func(t *testing.T, db *inMemory.DB, testName string) {
				for _, data := range testData[testName] {
					_, err := db.GetUrl(context.Background(), data.alias)
					assert.ErrorIs(t, err, data.expectedErr)
				}
			},
		},
	}

	for _, test := range tests {
		db := inMemory.NewDB()
		test.prepare(db)
		t.Run(test.name, func(t *testing.T) {
			test.verifyResult(t, db, test.name)
		})
	}
}
