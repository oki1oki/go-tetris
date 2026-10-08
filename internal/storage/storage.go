package storage

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func LoadHighScore() int {
	path, err := filePath()
	if err != nil {
		log.Println("storage: config dir:", err)
		return 0
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Println("storage: read:", err)
		}
		return 0
	}

	score, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		log.Println("storage: parse:", err)
		return 0
	}

	return score
}

func SaveHighScore(score int) error {
	path, err := filePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(strconv.Itoa(score)), 0o644)
}

func filePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "tetris", "highscore.txt"), nil
}
