package audio

import (
	"os"
	"time"
	"voz-em-texto/internal/system"
)

var UltimoAudioGerado string

func PrepararPastas() {

	pastas := []string{
		system.AudioDir(),
		system.InputDir(),
		system.OutputDir(),
	}

	for _, pasta := range pastas {
		if _, err := os.Stat(pasta); os.IsNotExist(err) {
			os.MkdirAll(pasta, 0755)
		}
	}
}

func CopiarArquivo(origem, destino string) error {
	input, err := os.ReadFile(origem)
	if err != nil {
		return err
	}
	return os.WriteFile(destino, input, 0644)
}

func Timestamp() string {
	return time.Now().Format("02012006_150405")
}
