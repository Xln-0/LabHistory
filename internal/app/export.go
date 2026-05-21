package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (m *Model) ExportEnv() error {

	dir := filepath.Join(os.Getenv("HOME"), ".config/labhistory")
	file := filepath.Join(dir, "env.sh")

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	lines, err := readEnvFile(file)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// -------------------------
	// HOST MODE
	// -------------------------
	if m.focus == 0 {

		if len(m.hosts) == 0 || m.hostSelected >= len(m.hosts) {
			return nil
		}

		host := m.hosts[m.hostSelected]

		lines = setEnvVar(lines, "HOST_NAME", host.Name)
		lines = setEnvVar(lines, "HOST_IP", host.IP)
		lines = setEnvVar(lines, "HOST_DOMAIN", host.Domain)
		lines = setEnvVar(lines, "HOST_ROLE", host.Role)
	}

	// -------------------------
	// USER MODE
	// -------------------------
	if m.focus == 1 {

		if len(m.users) == 0 || m.userSelected >= len(m.users) {
			return nil
		}

		user := m.users[m.userSelected]

		lines = setEnvVar(lines, "USER_NAME", user.Username)
		lines = setEnvVar(lines, "USER_PASSWD", user.Password)
		lines = setEnvVar(lines, "USER_HASH", user.Hash)
	}

	return writeEnvFile(file, lines)
}

// -------------------------
// CORE LOGIC (sed-like)
// -------------------------

func setEnvVar(lines []string, key, value string) []string {

	prefix := "export " + key + "="

	replaced := false

	for i, l := range lines {

		l = strings.TrimSpace(l)

		if strings.HasPrefix(l, prefix) {
			lines[i] = fmt.Sprintf("export %s='%s'", key, value)
			replaced = true
		}
	}

	if !replaced {
		lines = append(lines, fmt.Sprintf("export %s='%s'", key, value))
	}

	return lines
}

// -------------------------
// FILE IO
// -------------------------

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
