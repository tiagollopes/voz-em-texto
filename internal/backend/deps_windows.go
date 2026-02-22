//go:build windows

package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func ChecarDependencias() error {
	fmt.Println("⚙️ Verificando dependências internas...")

	// 1. Pega o caminho de onde o programa está executando
	execPath, _ := os.Executable()
	baseDir := filepath.Dir(execPath)

	// 2. Monta o caminho para a nossa pasta bin
	caminhoLocal := filepath.Join(baseDir, "bin", runtime.GOOS, "ffmpeg.exe")

	// 3. Verifica se o ffmpeg.exe está lá
	if _, err := os.Stat(caminhoLocal); err == nil {
		fmt.Println("✅ FFmpeg interno encontrado em:", caminhoLocal)
		return nil
	}

	fmt.Println("⚠️ FFmpeg interno não encontrado. Certifique-se de que a pasta 'bin' está junto com o executável.")
	return fmt.Errorf("ffmpeg.exe ausente na pasta bin/%s", runtime.GOOS)
}
