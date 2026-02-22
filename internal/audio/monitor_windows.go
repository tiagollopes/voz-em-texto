//go:build windows

package audio

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type MonitorInfo struct {
	Nome   string
	Driver string // "dshow" ou "wasapi_custom"
}

func DetectarMonitor() (MonitorInfo, error) {
	exePath, _ := os.Executable()
	baseDir := filepath.Dir(exePath)
	ffmpegPath := filepath.Join(baseDir, "bin", "windows", "ffmpeg.exe")

	// 1️⃣ Tenta detectar via DSHOW (Stereo Mix)
	cmd := exec.Command(ffmpegPath, "-list_devices", "true", "-f", "dshow", "-i", "dummy")

	// 2. ADICIONE ESTA LINHA AQUI:
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Run()

	scanner := bufio.NewScanner(&stderr)
	for scanner.Scan() {
		linha := scanner.Text()
		if strings.Contains(linha, "Mixagem estéreo") ||
		   strings.Contains(linha, "Stereo Mix") ||
		   strings.Contains(linha, "What U Hear") {

			start := strings.Index(linha, "\"")
			end := strings.LastIndex(linha, "\"")
			if start != -1 && end != -1 && start != end {
				return MonitorInfo{
					Nome:   linha[start+1 : end],
					Driver: "dshow",
				}, nil
			}
		}
	}

	// 2️⃣ Fallback para seu binário customizado WASAPI
	return MonitorInfo{
		Nome:   "System Loopback (WASAPI)",
		Driver: "wasapi_custom",
	}, nil
}
