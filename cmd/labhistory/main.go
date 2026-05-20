package main

import (
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Xln-0/labhistory/internal/db"
	"github.com/Xln-0/labhistory/internal/tui"
)

func main() {
	path, err := getDBPath()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Init(database); err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(
		tui.New(database),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func getDBPath() (string, error) {

	cfg, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(cfg, "labhistory")

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(dir, "labhistory.db"), nil
}
