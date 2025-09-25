package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// RunCommand exécute une commande système et retourne sa sortie
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing '%s %v': %w", name, args, err)
	}
	return strings.TrimSpace(string(output)), nil
}

// FormatSpeed convertit des bytes/seconde en Mb/s formaté
func FormatSpeed(bytesPerSecond float64) string {
	bitsPerSecond := bytesPerSecond * 8
	megaBitsPerSecond := bitsPerSecond / 1000000
	
	return fmt.Sprintf("%.1f Mb/s", megaBitsPerSecond)
}

// CountLines compte le nombre de lignes dans une chaîne
func CountLines(output string) int {
	if output == "" {
		return 0
	}
	return len(strings.Split(output, "\n"))
}