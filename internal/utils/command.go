package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// RunCommand executes a system command and returns its output
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing '%s %v': %w", name, args, err)
	}
	return strings.TrimSpace(string(output)), nil
}

// FormatSpeed converts bytes/second to formatted Mb/s
func FormatSpeed(bytesPerSecond float64) string {
	bitsPerSecond := bytesPerSecond * 8
	megaBitsPerSecond := bitsPerSecond / 1000000
	
	return fmt.Sprintf("%.1f Mb/s", megaBitsPerSecond)
}

// CountLines counts the number of lines in a string
func CountLines(output string) int {
	if output == "" {
		return 0
	}
	return len(strings.Split(output, "\n"))
}