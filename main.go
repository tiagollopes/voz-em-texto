//sudo apt update
//sudo apt install -y cmake build-essential

/*
git clone https://github.com/tiagollopes/voz-em-texto
cd voz-em-texto
chmod +x install.sh
./install.sh
go run main.go
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var ultimoAudioGerado string

// ===============================
// Detectar monitor
// ===============================
func detectarMonitor() (string, error) {

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

func prepararPastas() {

	pastas := []string{
		"audio",
		"input",
		"output",
	}

	for _, pasta := range pastas {

		if _, err := os.Stat(pasta); os.IsNotExist(err) {
			os.Mkdir(pasta, 0755)
		}
	}
}

// ===============================
// Gravar áudio
// ===============================
func gravarAudio(monitor string) error {

	fmt.Println("♪ Gravando áudio...")

	cmd := exec.Command(
		"ffmpeg",
		"-y", // força sobrescrever
		"-f", "pulse",
		"-i", monitor,
		"-t", "30",
		"-ac", "1",
		"-ar", "16000",
		"-b:a", "64k",
		"audio/audio.mp3",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ===============================
// Instalacao Whisper
// ===============================
func instalarWhisper() error {

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
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	fmt.Println("⚙ Compilando Whisper...")

	cmd := exec.Command("make")
	cmd.Dir = "./whisper"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("⬇️ Baixando modelo...")

	cmd = exec.Command(
		"bash",
		"./models/download-ggml-model.sh",
		"base",
	)
	cmd.Dir = "./whisper"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ===============================
// Transcrever
// ===============================
func transcrever() error {

	duracao, _ := duracaoAudio()

	cmd := exec.Command(
		"./whisper/build/bin/whisper-cli",
		"-m", "whisper/models/ggml-base.bin",
		"-f", "audio/audio.mp3",
		"-l", "pt", // força português
		"-otxt",
	)

	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Start()
	if err != nil {
		return err
	}

	stopSpinner := make(chan bool)
	go spinnerPercent(stopSpinner, duracao)

	err = cmd.Wait()

	stopSpinner <- true

	fmt.Println("\n✅ Transcrição finalizada!")

	return err
}

// ===============================
// Checagem Dependencia Linux
// ===============================
func checarDependencias() error {

	fmt.Println("⚙ Verificando dependências...")

	deps := []string{
		"ffmpeg",
		"cmake",
		"make",
		"gcc",
		"git",
	}

	for _, dep := range deps {

		_, err := exec.LookPath(dep)
		if err != nil {

			fmt.Printf("❌ Dependência faltando: %s\n", dep)
			fmt.Println("")
			fmt.Println("➡ Rode primeiro:")
			fmt.Println("./install.sh")
			fmt.Println("")

			return err
		}
	}

	fmt.Println("✅ Dependências OK")
	return nil
}

func copiarArquivo(origem, destino string) error {

	input, err := os.ReadFile(origem)
	if err != nil {
		return err
	}

	return os.WriteFile(destino, input, 0644)
}

func gravarAteParar(monitor string) error {

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
	go spinnerGravando(stopSpinner, inicio)

	fmt.Scanln()

	cmd.Process.Signal(os.Interrupt)
	cmd.Wait()

	stopSpinner <- true

	fmt.Println("\n✅ Gravação parada.")
	ts := timestamp()

	nomeBase := "gravacao_" + ts

	origem := "audio/audio.mp3"
	destino := "output/" + nomeBase + ".mp3"

	copiarArquivo(origem, destino)

	// guarda na sessão
	ultimoAudioGerado = nomeBase

	fmt.Println("✅ Áudio salvo em:", destino)

	return nil
}

func spinnerGravando(stopChan chan bool, inicio time.Time) {

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

func duracaoAudio() (float64, error) {

	cmd := exec.Command(
		"ffprobe",
		"-i", "audio/audio.mp3",
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

func spinnerPercent(stopChan chan bool, duracao float64) {

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

func piscarTexto(stopChan chan bool) {

	for {
		select {
		case <-stopChan:

			// Limpa linha de baixo
			fmt.Print("\033[2K\r")
			return

		default:

			// Vai pra linha de baixo
			fmt.Print("\033[s")  // salva cursor
			fmt.Print("\033[1B") // desce 1 linha

			fmt.Print("\rPressione ENTER para parar gravação")
			time.Sleep(500 * time.Millisecond)

			fmt.Print("\r                                   ")
			time.Sleep(500 * time.Millisecond)

			fmt.Print("\033[u") // volta cursor
		}
	}
}

func menu() string {

	var opcao string

	fmt.Println("")
	fmt.Println("1 - Gravar áudio")
	fmt.Println("2 - Transcrever áudio existente")
	fmt.Println("0 - Sair")
	fmt.Print("Escolha: ")

	fmt.Scanln(&opcao)

	return opcao
}

func transcreverArquivo() {

	arquivos, err := os.ReadDir("input")
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
	caminho := "input/" + nome

	fmt.Println("\n⚙ Transcrevendo:", nome)

	// Instala whisper se precisar
	err = instalarWhisper()
	if err != nil {
		fmt.Println("Erro instalando Whisper.")
		return
	}

	// Duração do áudio
	duracao, _ := duracaoArquivo(caminho)

	cmd := exec.Command(
		"./whisper/build/bin/whisper-cli",
		"-m", "whisper/models/ggml-base.bin",
		"-f", caminho,
		"-l", "pt",
		"-otxt",
	)

	cmd.Stdout = nil
	cmd.Stderr = nil

	err = cmd.Start()
	if err != nil {
		fmt.Println("Erro iniciando transcrição.")
		return
	}

	// Spinner percentual
	stopSpinner := make(chan bool)
	go spinnerPercent(stopSpinner, duracao)

	err = cmd.Wait()

	stopSpinner <- true

	if err != nil {
		fmt.Println("❌ Erro na transcrição.")
		return
	}

	// Move saída
	os.Rename(
		caminho+".txt",
		"output/"+nome+".txt",
	)

	fmt.Println("\n✅ Transcrição salva em output/")
}

func duracaoArquivo(caminho string) (float64, error) {

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

func timestamp() string {
	return time.Now().Format("02012006_150405")
}

func transcreverUltimo() error {

	if ultimoAudioGerado == "" {
		fmt.Println("Nenhuma gravação na sessão.")
		return fmt.Errorf("nenhuma gravação na sessão")
	}

	audioPath := "output/" + ultimoAudioGerado + ".mp3"

	fmt.Println("⚙ Transcrevendo:", ultimoAudioGerado)

	duracao, _ := duracaoArquivo(audioPath)

	cmd := exec.Command(
		"./whisper/build/bin/whisper-cli",
		"-m", "whisper/models/ggml-base.bin",
		"-f", audioPath,
		"-l", "pt",
		"-otxt",
	)

	cmd.Stdout = nil
	cmd.Stderr = nil

	// Inicia processo
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Spinner enquanto roda
	stopSpinner := make(chan bool)
	go spinnerPercent(stopSpinner, duracao)

	// Aguarda finalizar
	err = cmd.Wait()

	// Para spinner
	stopSpinner <- true

	if err != nil {
		return err
	}

	txtOrigem := audioPath + ".txt"
	txtDestino := "output/" + ultimoAudioGerado + ".txt"

	os.Rename(txtOrigem, txtDestino)

	fmt.Println(" Texto salvo em:", txtDestino)

	return nil
}

// ===============================
// MAIN PIPELINE
// ===============================
func main() {

	err := checarDependencias()
	if err != nil {
		return
	}

	prepararPastas()

	opcao := menu()

	switch opcao {

	// =========================
	// 1 - GRAVAR
	// =========================
	case "1":

		monitor, err := detectarMonitor()
		if err != nil {
			fmt.Println("❌ Erro:", err)
			return
		}

		err = gravarAteParar(monitor)
		if err != nil {
			fmt.Println("❌ Erro na gravação:", err)
			return
		}

		fmt.Println("✅ Gravação concluída!")

		// Pergunta se quer transcrever
		var resp string
		fmt.Print("Deseja transcrever o áudio? (S/N): ")
		fmt.Scanln(&resp)

		resp = strings.ToLower(strings.TrimSpace(resp))

		if resp == "s" {

			err = instalarWhisper()
			if err != nil {
				fmt.Println("❌ Erro instalando Whisper:", err)
				return
			}

			err = transcreverUltimo() //err = transcrever()
			if err != nil {
				fmt.Println("❌ Erro na transcrição:", err)
				return
			}

			/*os.Rename(
				"audio/audio.mp3.txt",
				"output/transcricao.txt",
			)
			*/
			fmt.Println("✅ Transcrição salva em transcricao.txt")
		}

		fmt.Println("✅  Processo finalizado.")

	// =========================
	// 2 - TRANSCRIBIR ARQUIVO
	// =========================
	case "2":

		transcreverArquivo()

	// =========================
	// 0 - SAIR
	// =========================
	case "0":

		fmt.Println("Encerrado.")
		return

	default:

		fmt.Println("Opção inválida.")
	}
}
