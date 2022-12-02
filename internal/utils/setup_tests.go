package utils

import (
	"testing"

	"github.com/iryzzh/y-gophkeeper/internal/models"
	"github.com/iryzzh/y-gophkeeper/internal/rand"

	"github.com/iryzzh/y-gophkeeper/internal/config"
)

// TestConfig returns the test configuration.
func TestConfig(tb testing.TB) (*config.ServerCfg, error) {
	tb.Helper()

	cfg, err := config.NewServerConfig()
	if err != nil {
		return nil, err
	}
	cfg.DB.DSN = ":memory:"

	cfg.DB.MigrationsPath = "./migrations"

	return cfg, nil
}

// TestItem returns a randomly generated item.
func TestItem(t *testing.T, userID string) *models.Item {
	t.Helper()
	return &models.Item{
		UserID: userID,
		Meta:   rand.String(10), //nolint:gomnd
		ItemData: &models.ItemData{
			Data: []byte(rand.String(10)), //nolint:gomnd
		},
	}
}
