//go:build !windows

package transcribe

import (
	"fmt"
	"os"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"voz-em-texto/internal/system"
	"voz-em-texto/internal/audio"
	"voz-em-texto/internal/progress"
	"path/filepath"
)


var TranscribeCmd *exec.Cmd

// ===============================
// Instalacao Whisper
// ===============================
func InstalarWhisper() error {

	fmt.Println("📦 Verificando Whisper.cpp...")

	// Se já existe binário → pula
	if _, err := os.Stat("./whisper/build/bin/whisper-cli"); err == nil {
		fmt.Println("✅ Whisper já instalado.")
		return nil
	}

	// Se pasta existe mas não compilado
	if _, err := os.Stat("./whisper"); err == nil {
		fmt.Println("✅ Pasta whisper já existe. Pulando clone...")
	} else {

		fmt.Println("⬇️ Baixando Whisper.cpp...")

		cmd := exec.Command(
			"git", "clone",
			"https://github.com/ggerganov/whisper.cpp",
			"whisper",
		)
		/*cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr*/
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	fmt.Println("⚙ Compilando Whisper...")

	cmd := exec.Command("make")
	cmd.Dir = "./whisper"
	/*cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr*/
	cmd.Stdout = io.Discard
        cmd.Stderr = io.Discard

	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("⬇️ Baixando modelo...")

	cmd = exec.Command(
		"bash",
		"./models/download-ggml-model.sh",
		//"base",
		"tiny",
	)
	cmd.Dir = "./whisper"
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	return cmd.Run()
}

func DuracaoArquivo(caminho string) (float64, error) {

	cmd := exec.Command(
		"ffprobe",
		"-i", caminho,
		"-show_entries", "format=duration",
		"-v", "quiet",
		"-of", "csv=p=0",
	)

	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	duracao, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)

	return duracao, nil
}

func TranscreverUltimo() error {

	if audio.UltimoAudioGerado == "" {
		fmt.Println("Nenhuma gravação na sessão.")
		return fmt.Errorf("nenhuma gravação na sessão")
	}

	audioPath := system.OutputDir() + "/" + audio.UltimoAudioGerado + ".mp3"

	fmt.Println("⚙ Transcrevendo:", audio.UltimoAudioGerado)

	duracao, _ := DuracaoArquivo(audioPath)
	binario := system.WhisperBinary()
	modelo := system.WhisperDir() + "/models/ggml-tiny.bin"

	TranscribeCmd = exec.Command(
		binario,
		//"-m", "whisper/models/ggml-base.bin",
		"-m", modelo,
		"-f", audioPath,
		"-l", "pt",
		"-otxt",
	)

	TranscribeCmd.Stdout = nil
	TranscribeCmd.Stderr = nil

	// Inicia processo
	err := TranscribeCmd.Start()
	if err != nil {
		return err
	}

	// Spinner enquanto roda
	stopSpinner := make(chan bool)
	go progress.SpinnerPercent(stopSpinner, duracao)

	// Aguarda finalizar
	err = TranscribeCmd.Wait()


	// Para spinner
	stopSpinner <- true

	if err != nil {

	// Se foi cancelado manualmente
	if strings.Contains(err.Error(), "signal") ||
		strings.Contains(err.Error(), "killed") {

		fmt.Println("\n⛔ Transcrição cancelada.")
		return nil
	}

	return err
	}

	txtOrigem := audioPath + ".txt"
	txtDestino := system.OutputDir() + "/" + audio.UltimoAudioGerado + ".txt"

	os.Rename(txtOrigem, txtDestino)

	fmt.Println(" Texto salvo em:", txtDestino)

	return nil
}
func TranscreverCaminho(caminho string) error {

	fmt.Println("⚙ Transcrevendo:", caminho)

	// Garante Whisper instalado
	err := InstalarWhisper()
	if err != nil {
		return err
	}

	// Duração do áudio
	duracao, _ := DuracaoArquivo(caminho)

	binario := system.WhisperBinary()
	modelo := system.WhisperDir() + "/models/ggml-tiny.bin"

	TranscribeCmd = exec.Command(
		binario,
		"-m", modelo,
		"-f", caminho,
		"-l", "pt",
		"-otxt",
	)

	TranscribeCmd.Stdout = nil
	TranscribeCmd.Stderr = nil

	// Inicia processo
	err = TranscribeCmd.Start()
	if err != nil {
		return err
	}

	// Spinner
	stopSpinner := make(chan bool)
	go progress.SpinnerPercent(stopSpinner, duracao)

	// Aguarda
	err = TranscribeCmd.Wait()

	stopSpinner <- true

	if err != nil {

		if strings.Contains(err.Error(), "signal") ||
			strings.Contains(err.Error(), "killed") {

			fmt.Println("\n⛔ Transcrição cancelada.")
			return nil
		}

		return err
	}

	// Move TXT
	nome := filepath.Base(caminho)

	os.Rename(
		caminho+".txt",
		system.OutputDir()+"/"+nome+".txt",
	)

	fmt.Println("✅ Transcrição salva em output/")

	return nil
}

func TranscreverArquivo() {

	arquivos, err := os.ReadDir(system.InputDir())

	if err != nil {
		fmt.Println("❌ Erro lendo pasta input.")
		return
	}

	if len(arquivos) == 0 {
		fmt.Println("⚠ Nenhum áudio encontrado em /input")
		return
	}

	fmt.Println("\n==============================")
	fmt.Println(" Arquivos encontrados")
	fmt.Println("==============================\n")

	for i, arq := range arquivos {
		fmt.Printf("[%d] %s\n", i+1, arq.Name())
	}

	var escolha int
	fmt.Print("\nEscolha o arquivo: ")
	fmt.Scanln(&escolha)

	if escolha < 1 || escolha > len(arquivos) {
		fmt.Println("❌ Opção inválida.")
		return
	}

	nome := arquivos[escolha-1].Name()
	caminho := system.InputDir() + "/" + nome

	fmt.Println("\n⚙ Transcrevendo:", nome)

	// Instala whisper se precisar
	err = InstalarWhisper()
	if err != nil {
		fmt.Println("Erro instalando Whisper.")
		return
	}

	// Duração do áudio
	duracao, _ := DuracaoArquivo(caminho)
	binario := system.WhisperBinary()
	modelo := system.WhisperDir() + "/models/ggml-tiny.bin"

	TranscribeCmd = exec.Command(
		binario,
		//"-m", "whisper/models/ggml-base.bin",
		"-m", modelo,
		"-f", caminho,
		"-l", "pt",
		"-otxt",
	)

	TranscribeCmd.Stdout = nil
	TranscribeCmd.Stderr = nil

	err = TranscribeCmd.Start()
	if err != nil {
		fmt.Println("Erro iniciando transcrição.")
		return
	}

	// Spinner percentual
	stopSpinner := make(chan bool)
	go progress.SpinnerPercent(stopSpinner, duracao)

	err = TranscribeCmd.Wait()

	stopSpinner <- true


	if err != nil {

	// Se foi cancelado manualmente
	if strings.Contains(err.Error(), "signal") ||
		strings.Contains(err.Error(), "killed") {

		fmt.Println("\n⛔ Transcrição cancelada.")
		return
	}

	fmt.Println("❌ Erro na transcrição.")
	return
	}


	// Move saída
	os.Rename(
		caminho+".txt",
		system.OutputDir()+"/"+nome+".txt",
	)

	fmt.Println("\n✅ Transcrição salva em output/")
}

// ===============================
// Parar Transcrição
// ===============================
func PararTranscricao() error {

	if TranscribeCmd == nil {
		fmt.Println("⛔ Nenhuma transcrição ativa.")
		return nil
	}

	if TranscribeCmd.Process == nil {
		fmt.Println("⛔ Processo já finalizado.")
		return nil
	}

	// Mata processo principal
	_ = TranscribeCmd.Process.Kill()

	// Mata filhos whisper
	exec.Command("pkill", "-f", "whisper-cli").Run()

	fmt.Println("\n⛔ Transcrição cancelada pelo usuário.")

	return nil
}
