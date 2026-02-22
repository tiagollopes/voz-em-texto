package cli

import (
	"fmt"
	"time"

	"voz-em-texto/internal/audio"
)

func GravarInterativo() error {

	fmt.Println("⚙ Detectando monitor...")

	monitor, err := audio.DetectarMonitor()
	if err != nil {
		return err
	}

	err = audio.IniciarGravacao(monitor)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Pressione ENTER para parar")
	fmt.Println()

	inicio := time.Now()

	stopSpinner := make(chan bool)
	go spinnerGravando(stopSpinner, inicio)

	fmt.Scanln()

	stopSpinner <- true

	err = audio.PararGravacao()
	if err != nil {
		return err
	}

	fmt.Println("\n✅ Gravação finalizada")

	return nil
}
