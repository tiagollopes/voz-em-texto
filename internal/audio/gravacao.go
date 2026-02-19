package audio

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"os"
	"time"
	"voz-em-texto/internal/system"
)

var cmdGravacao *exec.Cmd
var UltimoAudioGerado string

// ===============================
// Detectar monitor
// ===============================
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

			//fmt.Println("✅ Monitor encontrado:", monitor)
			fmt.Println("✅ Monitor de áudio encontrado:")
			return monitor, nil
		}
	}

	return "", fmt.Errorf("nenhum monitor encontrado")
}

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

func IniciarGravacao(monitor string) error {
arquivo := system.AudioDir() + "/audio.mp3"
	cmdGravacao = exec.Command(
		"ffmpeg",
		"-y",
		"-loglevel", "quiet",
		"-f", "pulse",
		"-i", monitor,
		"-ac", "1",
		"-ar", "16000",
		"-b:a", "64k",
		arquivo,
	)

	return cmdGravacao.Start()
}
func PararGravacao() error {

	if cmdGravacao == nil {
		return fmt.Errorf("nenhuma gravação em andamento")
	}

	err := cmdGravacao.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	cmdGravacao.Wait()

	ts := Timestamp()
	nomeBase := "gravacao_" + ts

	origem := system.AudioDir() + "/audio.mp3"
	destino := system.OutputDir() + "/" + nomeBase + ".mp3"

	CopiarArquivo(origem, destino)

	UltimoAudioGerado = nomeBase

	return nil
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


///cli
func GravarAteParar(monitor string) error {

	// Espaço UI
	fmt.Println()
	fmt.Println()
	fmt.Println()

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-loglevel", "quiet",
		"-f", "pulse",
		"-i", monitor,
		"-ac", "1",
		"-ar", "16000",
		"-b:a", "64k",
		"audio/audio.mp3",
	)

	err := cmd.Start()
	if err != nil {
		return err
	}

	inicio := time.Now()

	stopSpinner := make(chan bool)
	go SpinnerGravando(stopSpinner, inicio)

	fmt.Scanln()

	cmd.Process.Signal(os.Interrupt)
	cmd.Wait()

	stopSpinner <- true

	fmt.Println("\n✅ Gravação parada.")
	ts := Timestamp()

	nomeBase := "gravacao_" + ts

	origem := system.AudioDir() + "/audio.mp3"
	destino := system.OutputDir() + "/" + nomeBase + ".mp3"

	CopiarArquivo(origem, destino)

	// guarda na sessão
	UltimoAudioGerado = nomeBase

	fmt.Println("✅ Áudio salvo em:", destino)

	return nil
}

func SpinnerGravando(stopChan chan bool, inicio time.Time) {

	chars := []string{"/", "-", "\\", "|"}
	i := 0

	for {
		select {
		case <-stopChan:

			// Limpa as 3 linhas
			fmt.Print("\r\033[K")
			fmt.Print("\033[1B\r\033[K")
			fmt.Print("\033[1B\r\033[K")
			fmt.Print("\033[2A")

			return

		default:

			decorrido := time.Since(inicio)

			min := int(decorrido.Minutes())
			seg := int(decorrido.Seconds()) % 60

			// ---------- Linha 1 ----------
			fmt.Printf(
				"\r● REC %02d:%02d %s",
				min,
				seg,
				chars[i],
			)

			// ---------- Barra fake ----------
			total := 10
			pos := (seg % total)

			barra := ""
			for j := 0; j < total; j++ {
				if j <= pos {
					barra += "█"
				} else {
					barra += "░"
				}
			}

			// ---------- Linha 2 ----------
			fmt.Print("\n" + barra)

			// ---------- Linha 3 ----------
			fmt.Print("\nPressione ENTER para parar")

			// Volta cursor pra linha 1
			fmt.Print("\033[2A")

			time.Sleep(200 * time.Millisecond)
			i = (i + 1) % len(chars)
		}
	}
}
