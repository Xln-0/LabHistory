package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/Xln-0/labhistory/internal/app/components"
	"github.com/Xln-0/labhistory/internal/db"
)

// -------------------------
// Set Environment Variables
// -------------------------

func setEnvVar(lines []string, key, value string) []string {

	prefix := "export " + key + "="

	for i := range lines {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), prefix) {
			lines[i] = fmt.Sprintf("export %s='%s'", key, value)
			return lines
		}
	}

	return append(lines, fmt.Sprintf("export %s='%s'", key, value))
}

func readEnvFile(file string) ([]string, error) {

	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	// remove trailing empty line if any
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines, nil
}

func writeEnvFile(file string, lines []string) error {

	content := strings.Join(lines, "\n") + "\n"

	return os.WriteFile(file, []byte(content), 0644)
}

func CapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func subdomainsToList(subs []db.Subdomain) []components.ListItem {

	items := make([]components.ListItem, len(subs))

	for i, s := range subs {

		items[i] = components.ListItem{
			ID:   s.ID,
			Name: s.Subdomain,
		}
	}

	return items
}
