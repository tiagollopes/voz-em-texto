package cli

import (
	"fmt"
	"time"
)

func spinnerGravando(stopChan chan bool, inicio time.Time) {

	chars := []string{"/", "-", "\\", "|"}
	i := 0

	for {
		select {

		case <-stopChan:

			// limpa linhas
			fmt.Print("\r\033[K")
			fmt.Print("\033[1B\r\033[K")
			fmt.Print("\033[1B\r\033[K")
			fmt.Print("\033[2A")

			return

		default:

			decorrido := time.Since(inicio)

			min := int(decorrido.Minutes())
			seg := int(decorrido.Seconds()) % 60

			fmt.Printf(
				"\r● REC %02d:%02d %s",
				min,
				seg,
				chars[i],
			)

			time.Sleep(200 * time.Millisecond)
			i = (i + 1) % len(chars)
		}
	}
}
