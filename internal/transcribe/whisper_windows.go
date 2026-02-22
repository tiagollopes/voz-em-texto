//go:build windows

package transcribe

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"voz-em-texto/internal/audio"
	"voz-em-texto/internal/progress"
	"strconv"
	"syscall"
)

var TranscribeCmd *exec.Cmd

// TranscreverUltimo é a função chamada pela OPÇÃO 2
func TranscreverUltimo() error {
	arquivo := audio.UltimoAudioGerado
	if arquivo == "" {
		return fmt.Errorf("nenhuma gravação encontrada nesta sessão")
	}

	// AJUSTE 1: Garante que o caminho seja absoluto e tenha extensão
	if !strings.HasSuffix(strings.ToLower(arquivo), ".wav") {
		arquivo += ".wav"
	}

	// AJUSTE 2: Verifica se o arquivo REALMENTE existe no disco antes de tentar
	caminhoAbs, _ := filepath.Abs(arquivo)
	if _, err := os.Stat(caminhoAbs); os.IsNotExist(err) {
		// Tenta procurar na pasta output se não achar na raiz
		exePath, _ := os.Executable()
		caminhoNaOutput := filepath.Join(filepath.Dir(exePath), "output", filepath.Base(arquivo))
		if _, err := os.Stat(caminhoNaOutput); err == nil {
			caminhoAbs = caminhoNaOutput
		} else {
			return fmt.Errorf("arquivo não encontrado: %s", filepath.Base(arquivo))
		}
	}

	fmt.Printf("✅ Iniciando transcrição de: %s\n", filepath.Base(caminhoAbs))
	return TranscreverCaminho(caminhoAbs)
}

// TranscreverArquivo é a função chamada pela OPÇÃO 3
func TranscreverArquivo() {
	exePath, _ := os.Executable()
	baseDir := filepath.Dir(exePath)
	pastaInput := filepath.Join(baseDir, "input")

	arquivos, err := os.ReadDir(pastaInput)
	if err != nil || len(arquivos) == 0 {
		fmt.Println("❌ Nenhuma gravação encontrada na pasta input.")
		return
	}

	fmt.Println("\n--- Arquivos .wav disponíveis ---")
	var listaWav []string
	count := 1
	for _, f := range arquivos {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".wav") {
			fmt.Printf("%d - %s\n", count, f.Name())
			listaWav = append(listaWav, f.Name())
			count++
		}
	}

	if len(listaWav) == 0 {
		fmt.Println("❌ Nenhum arquivo .wav encontrado.")
		return
	}

	fmt.Print("\nEscolha o número do arquivo (ou 0 para cancelar): ")
	var escolha int
	fmt.Scanln(&escolha)

	if escolha == 0 { return }

	if escolha < 1 || escolha > len(listaWav) {
		fmt.Println("❌ Escolha inválida.")
		return
	}

	caminhoCompleto := filepath.Join(pastaInput, listaWav[escolha-1])
	TranscreverCaminho(caminhoCompleto)
}

// 1. AJUSTE: Adicionado "error" no final da assinatura
func TranscreverCaminho(caminho string) error {
	exePath, _ := os.Executable()
	baseDir := filepath.Dir(exePath)

	caminhoAbs, _ := filepath.Abs(caminho)
	binDir := filepath.Join(baseDir, "bin", "windows")
	binario := filepath.Join(binDir, "whisper-cli.exe")
	modeloAbs, _ := filepath.Abs(filepath.Join(baseDir, "bin", "models", "ggml-tiny.bin"))

	if _, err := os.Stat(caminhoAbs); os.IsNotExist(err) {
		fmt.Printf("❌ Arquivo não encontrado: %s\n", caminhoAbs)
		// 2. AJUSTE: Agora retorna o erro em vez de vazio
		return fmt.Errorf("arquivo não encontrado")
	}

	duracao, errDur := DuracaoArquivo(caminhoAbs)
	if errDur != nil || duracao == 0 {
		duracao = 30
	}

	//nomeSemExt := strings.TrimSuffix(caminhoAbs, filepath.Ext(caminhoAbs))
	nomeBase := strings.TrimSuffix(filepath.Base(caminhoAbs), filepath.Ext(caminhoAbs))

	TranscribeCmd = exec.Command(binario,
	"-m", modeloAbs,
	"-f", caminhoAbs,
	"-l", "pt",
	"-otxt",
	"-of", nomeBase,
	)



	TranscribeCmd.Dir = binDir
	TranscribeCmd.Stdout = nil
	TranscribeCmd.Stderr = nil
	TranscribeCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	err := TranscribeCmd.Start()
	if err != nil {
		fmt.Printf("❌ Falha ao iniciar o Whisper: %v\n", err)
		// 3. AJUSTE: Retorna o erro da inicialização
		return err
	}

	stopSpinner := make(chan bool)
	go progress.SpinnerPercent(stopSpinner, duracao)

	err = TranscribeCmd.Wait()
	stopSpinner <- true

	if err != nil {
		fmt.Printf("\n❌ O Whisper encerrou com erro: %v\n", err)
		// 4. AJUSTE: Retorna o erro do processo
		return err
	}

	txtGerado := filepath.Join(binDir, nomeBase+".txt")

	if _, err := os.Stat(txtGerado); err != nil {
	fmt.Printf("\n❌ Erro: O arquivo .txt não foi gerado em: %s\n", txtGerado)
	return fmt.Errorf("arquivo txt não gerado")
	}

	pastaOutput := filepath.Join(baseDir, "output")
	os.MkdirAll(pastaOutput, 0755)

	destinoFinal := filepath.Join(pastaOutput, nomeBase+".txt")

	// Remove se já existir
	_ = os.Remove(destinoFinal)

	err = os.Rename(txtGerado, destinoFinal)
	if err != nil {
	return fmt.Errorf("erro ao mover txt para output: %v", err)
	}

	fmt.Printf("\n✅ Transcrição concluída com sucesso!")
	fmt.Printf("\n📂 Salvo em: %s\n", destinoFinal)

	// 6. AJUSTE: Retorna nil (sucesso) no final
	return nil
}

// InstalarWhisper apenas valida se as peças estão no lugar
func InstalarWhisper() error {
	exePath, _ := os.Executable()
	baseDir := filepath.Dir(exePath)
	binario := filepath.Join(baseDir, "bin", "windows", "whisper-cli.exe")

	if _, err := os.Stat(binario); err != nil {
		return fmt.Errorf("whisper-cli.exe não encontrado em bin/windows/")
	}
	return nil
}

func DuracaoArquivo(caminho string) (float64, error) {
	// AJUSTE: Buscar o ffprobe na pasta bin para garantir que funcione
	exePath, _ := os.Executable()
	ffprobePath := filepath.Join(filepath.Dir(exePath), "bin", "windows", "ffprobe.exe")

	cmd := exec.Command(ffprobePath,
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		caminho,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	f, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return 0, err
	}

	return f, nil
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

	// 1. Mata o processo principal (o objeto Cmd do Go)
	_ = TranscribeCmd.Process.Kill()

	// 2. Mata os processos filhos no Windows de forma invisível
	// /F = Forçar /T = Terminar processos filhos /IM = Nome da imagem (executável)
	cmdKill := exec.Command("taskkill", "/F", "/T", "/IM", "whisper-cli.exe")

	// Esconde a janela do terminal do taskkill
	cmdKill.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	_ = cmdKill.Run()

	fmt.Println("\n⛔ Transcrição cancelada pelo usuário.")

	return nil
}
