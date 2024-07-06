package initializer

import (
	"log/slog"
	"os"

	"github.com/illiakornyk/spy-cat/internal/storage/sqlite"
)

func InitializeStorage(storagePath string, logger *slog.Logger) *sqlite.Storage {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		logger.Error("Failed to open SQLite database", slog.Any("error", err))
		os.Exit(1)
	}
	return storage
}
