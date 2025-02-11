package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/VadimShara/url_shortening_service/internal/repo/postgres"
	"github.com/VadimShara/url_shortening_service/pkg/errs"
)

func setupTestDB(t *testing.T) *postgres.DB {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor:      wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
		AlwaysPullImage: false,
	}

	postgresC, err := testcontainers.GenericContainer(context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	require.NoError(t, err)

	host, err := postgresC.Host(context.Background())
	require.NoError(t, err)
	port, err := postgresC.MappedPort(context.Background(), "5432")
	require.NoError(t, err)

	db := postgres.NewDB("test", "test", host, port.Int(), "testdb", "disable")
	t.Cleanup(func() {
		db.CLose()
		_ = postgresC.Terminate(context.Background())
	})

	return db
}

func TestSaveUrl(t *testing.T) {
	db := setupTestDB(t)
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
			{"https://example.com", "qwertyuiop", errs.ErrUrlExists},
			{"https://newsite.com", "abcdefghij", errs.ErrAliasExists},
		},
	}

	tests := []struct {
		name         string
		verifyResult func(t *testing.T, db *postgres.DB, testName string)
	}{
		{
			name: "valid",
			verifyResult: func(t *testing.T, db *postgres.DB, testName string) {
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
			verifyResult: func(t *testing.T, db *postgres.DB, testName string) {
				for _, data := range testData[testName] {
					_, err := db.SaveUrl(context.Background(), data.url, data.alias)
					assert.ErrorIs(t, err, data.err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.verifyResult(t, db, test.name)
		})
	}
}

func TestGetUrl(t *testing.T) {
	db := setupTestDB(t)
	testData := map[string][]struct {
		alias       string
		expectedUrl string
		expectedErr error
	}{
		"valid": {
			{"exmpl", "https://example.com", nil},
			{"ggl", "https://google.com", nil},
		},
		"invalid": {
			{"unknown", "", errs.ErrUrlNotFound},
		},
	}

	tests := []struct {
		name         string
		prepare      func(db *postgres.DB)
		verifyResult func(t *testing.T, db *postgres.DB, testName string)
	}{
		{
			name: "valid",
			prepare: func(db *postgres.DB) {
				db.SaveUrl(context.Background(), "https://example.com", "exmpl")
				db.SaveUrl(context.Background(), "https://google.com", "ggl")
			},
			verifyResult: func(t *testing.T, db *postgres.DB, testName string) {
				for _, data := range testData[testName] {
					url, err := db.GetUrl(context.Background(), data.alias)
					assert.NoError(t, err)
					assert.Equal(t, data.expectedUrl, url)
				}
			},
		},
		{
			name:    "invalid",
			prepare: func(db *postgres.DB) {},
			verifyResult: func(t *testing.T, db *postgres.DB, testName string) {
				for _, data := range testData[testName] {
					_, err := db.GetUrl(context.Background(), data.alias)
					assert.ErrorIs(t, err, data.expectedErr)
				}
			},
		},
	}

	for _, test := range tests {
		test.prepare(db)
		t.Run(test.name, func(t *testing.T) {
			test.verifyResult(t, db, test.name)
		})
	}
}
