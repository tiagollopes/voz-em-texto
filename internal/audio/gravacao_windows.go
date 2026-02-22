//go:build windows

package audio

import (
	"fmt"
	"time"
	"os"
	"os/exec"
	"path/filepath"
	"syscall" // 1. IMPORTANTE: Adicionado syscall
	"voz-em-texto/internal/system"
)

var (
	cmdAudio          *exec.Cmd
	cmdPipe           *exec.Cmd
	caminhoArquivoWav string
)

func IniciarGravacao(monitor MonitorInfo) error {
	exePath, _ := os.Executable()
	baseDir := filepath.Dir(exePath)
	binDir := filepath.Join(baseDir, "bin", "windows")

	ffmpegPath := filepath.Join(binDir, "ffmpeg.exe")
	wasapiPath := filepath.Join(binDir, "wasapi_capture.exe")

	fileName := fmt.Sprintf("audio_%s.wav", Timestamp())
	caminhoArquivoWav = filepath.Join(system.AudioDir(), fileName)
	os.MkdirAll(system.AudioDir(), 0755)

	// Criamos o atributo de esconder janela uma vez para reutilizar
	hideWindow := &syscall.SysProcAttr{HideWindow: true}

	if monitor.Driver == "dshow" {
		cmdAudio = exec.Command(ffmpegPath,
			"-f", "dshow",
			"-i", fmt.Sprintf("audio=%s", monitor.Nome),
			"-y", caminhoArquivoWav,
		)
		// 2. AJUSTE: Esconde janela do FFmpeg dshow
		cmdAudio.SysProcAttr = hideWindow
		return cmdAudio.Start()
	}

	cmdPipe = exec.Command(wasapiPath)
	// 3. AJUSTE: Esconde janela do capturador wasapi
	cmdPipe.SysProcAttr = hideWindow

	stdout, err := cmdPipe.StdoutPipe()
	if err != nil {
		return err
	}

	cmdAudio = exec.Command(ffmpegPath,
		"-f", "f32le",
		"-ar", "48000",
		"-ac", "2",
		"-i", "pipe:0",
		"-y", caminhoArquivoWav,
	)
	cmdAudio.Stdin = stdout
	// 4. AJUSTE: Esconde janela do FFmpeg pipe
	cmdAudio.SysProcAttr = hideWindow

	if err := cmdPipe.Start(); err != nil {
		return err
	}
	return cmdAudio.Start()
}

func PararGravacao() error {
	if cmdAudio != nil && cmdAudio.Process != nil {
		matarProcesso(cmdAudio)
		if cmdPipe != nil {
			matarProcesso(cmdPipe)
		}

		// AJUSTE: Aguarda 200ms para o Windows soltar o arquivo .wav
		time.Sleep(200 * time.Millisecond)

		nome := filepath.Base(caminhoArquivoWav)
		destino := filepath.Join(system.OutputDir(), nome)
		os.MkdirAll(system.OutputDir(), 0755)

		// Tenta mover para a pasta de transcrição
		err := os.Rename(caminhoArquivoWav, destino)
		if err == nil {
			UltimoAudioGerado = destino
		} else {
			// Se o Windows travar o Rename, usamos o caminho original mesmo
			// Assim a transcrição NÃO quebra
			UltimoAudioGerado = caminhoArquivoWav
		}
		return nil
	}
	return fmt.Errorf("nenhuma gravação ativa")
}

func matarProcesso(c *exec.Cmd) {
	// 5. AJUSTE: Esconde a janela do próprio Taskkill
	kill := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", c.Process.Pid))
	kill.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	kill.Run()
}
