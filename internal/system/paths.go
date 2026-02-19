package system

import (
	"os"
	"path/filepath"
)

// BasePath retorna o diretório do executável
func BasePath() string {
	exe, err := os.Executable()
	if err != nil {
		panic("Erro ao localizar executável")
	}
	return filepath.Dir(exe)
}

// Diretórios do sistema

func AudioDir() string {
	return filepath.Join(BasePath(), "audio")
}

func OutputDir() string {
	return filepath.Join(BasePath(), "output")
}

func InputDir() string {
	return filepath.Join(BasePath(), "input")
}

func WhisperDir() string {
	return filepath.Join(BasePath(), "whisper")
}

// Binário whisper
func WhisperBinary() string {
	if IsWindows() {
		return filepath.Join(WhisperDir(), "build", "bin", "whisper-cli.exe")
	}
	return filepath.Join(WhisperDir(), "build", "bin", "whisper-cli")
}

