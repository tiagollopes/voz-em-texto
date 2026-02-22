//go:build linux

package audio

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func DetectarMonitor() (string, error) {

	fmt.Println("⚙ Procurando monitor de áudio...")

	cmd := exec.Command("pactl", "list", "sources", "short")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		linha := scanner.Text()

		if strings.Contains(linha, ".monitor") {
			partes := strings.Fields(linha)
			monitor := partes[1]

			fmt.Println("✅ Monitor de áudio encontrado:")
			return monitor, nil
		}
	}

	return "", fmt.Errorf("nenhum monitor encontrado")
}
