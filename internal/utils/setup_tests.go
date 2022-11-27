package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/rand"

	"github.com/iryzzh/gophkeeper/internal/config"
)

const (
	migrationsDirName = "migrations"
)

// FindMigrationsDir finds the migrations' directory.
func FindMigrationsDir(tb testing.TB) (string, error) {
	tb.Helper()

	baseDir := "."
	for i := 0; i < 10; i++ {
		files, err := os.ReadDir(baseDir)
		if err != nil {
			return "", err
		}
		for _, f := range files {
			if f.Name() == migrationsDirName {
				return baseDir, nil
			}
		}
		baseDir = fmt.Sprintf("../%s", baseDir)
	}

	return "", fmt.Errorf("could not find migrations directory")
}

// TestConfig returns the test configuration.
func TestConfig(tb testing.TB) (*config.Config, error) {
	tb.Helper()

	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	cfg.DB.DSN = ":memory:"

	mDir, err := FindMigrationsDir(tb)
	if err != nil {
		return nil, err
	}

	cfg.DB.MigrationsPath = fmt.Sprintf("%v/%v", mDir, migrationsDirName)

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
