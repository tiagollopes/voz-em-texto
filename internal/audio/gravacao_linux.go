//go:build linux

package audio

import (
	"fmt"
	"os"
	"os/exec"
	"voz-em-texto/internal/system"
)

var cmdGravacao *exec.Cmd

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
