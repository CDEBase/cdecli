package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func findFile(dir string, fileNames []string) string {
	var content string
	for _, f := range fileNames {
		data, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			continue
		}
		content = string(data)
		break
	}

	return content
}

func readFile(dir string, fileName string) ([]byte, error) {
	return os.ReadFile(filepath.Join(dir, fileName))
}

func runCommand(dir string, command string, env string) (bool, error) {
	if command == "" {
		return false, nil
	}

	cmd := exec.Command("bash", "-c", command)
	cmd.Env = append(
		os.Environ(),
		env,
		fmt.Sprintf("PATH=%s:%s", filepath.Join(dir, "node_modules", ".bin"), os.Getenv("PATH")),
	)

	cmd.Dir = dir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	fmt.Fprintf(os.Stderr, "# running: %s\n", command)

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("command execution was failed: %s (see output above)", err)
	}

	fmt.Fprintln(os.Stderr)

	return true, nil
}
