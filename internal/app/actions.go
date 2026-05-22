package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

		if len(m.hosts.Items) == 0 || m.hosts.Selected >= len(m.hosts.Items) {
			return nil
		}

		host := m.hosts.Items[m.hosts.Selected]

		lines = setEnvVar(lines, "HOST_NAME", host.Name)
		lines = setEnvVar(lines, "HOST_IP", host.IP)
		lines = setEnvVar(lines, "HOST_DOMAIN", host.Domain)
		lines = setEnvVar(lines, "HOST_ROLE", host.Role)
	}

	// -------------------------
	// USER MODE
	// -------------------------
	if m.focus == 1 {

		if len(m.users.Items) == 0 || m.users.Selected >= len(m.users.Items) {
			return nil
		}

		user := m.users.Items[m.users.Selected]

		lines = setEnvVar(lines, "USER_NAME", user.Username)
		lines = setEnvVar(lines, "USER_PASSWD", user.Password)
		lines = setEnvVar(lines, "USER_HASH", user.Hash)
	}

	return writeEnvFile(file, lines)
}

func writeEtcHostsWithSudo(line string, password string) error {

	cmd := exec.Command(
		"sudo",
		"-S",
		"sh",
		"-c",
		fmt.Sprintf("echo '%s' >> /etc/hosts", line),
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, password+"\n")
	}()

	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("%v: %s", err, string(out))
	}

	return nil
}
