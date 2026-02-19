package progress

import (
	"fmt"
	"time"
)

func SpinnerPercent(stopChan chan bool, duracao float64) {

	chars := []string{"/", "-", "\\", "|"}
	i := 0

	tempoEstimado := duracao * 3.0 // fator mais realista
	inicio := time.Now()

	for {
		select {
		case <-stopChan:
			fmt.Print("\r                                       \r")
			return
		default:

			decorrido := time.Since(inicio).Seconds()
			percent := (decorrido / tempoEstimado) * 100

			// Limita em 95%
			if percent > 95 {
				percent = 95
			}

			msg := fmt.Sprintf(
				"\r⚙ Transcrevendo... %.0f%% %s",
				percent,
				chars[i],
			)

			// Se chegou perto do fim
			if percent >= 95 {
				msg = fmt.Sprintf(
					"\r⚙ Finalizando transcrição... %s",
					chars[i],
				)
			}

			fmt.Print(msg)

			time.Sleep(200 * time.Millisecond)
			i = (i + 1) % len(chars)
		}
	}
}
