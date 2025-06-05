package minotedb

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func BackupDatabase(backupPath string, dbFile string) (int64, string, error) {

	ts := time.Now().String()[:22]
	dbBakFile := filepath.Join(backupPath, fmt.Sprintf("minote.db-%s", ts))

	src, err := os.Open(dbFile)
	if err != nil {
		return 0, dbBakFile, err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(dbBakFile)
	if err != nil {
		return 0, dbBakFile, err
	}
	defer dst.Close()

	// Copy source to destination
	bytesCopied, err := io.Copy(dst, src)
	if err != nil {
		return 0, dbBakFile, err
	}
	return bytesCopied, dbBakFile, nil
}
