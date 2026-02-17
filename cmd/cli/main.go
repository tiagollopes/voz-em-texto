package main

import (
	"fmt"
	"os"
	"voz-em-texto/internal/backend" //  go.mod
)

func main() {
	// Inicialização básica do sistema
	backend.PrepararPastas()
	err := backend.ChecarDependencias()
	if err != nil {
		fmt.Printf("⚠️ Erro de dependências: %v\n", err)
		os.Exit(1)
	}

	for {
		fmt.Println("\n==============================")
		fmt.Println("      VOZ EM TEXTO - CLI")
		fmt.Println("==============================")
		fmt.Println("1 - Gravar Áudio (Sistema)")
		fmt.Println("2 - Transcrever Última Gravação")
		fmt.Println("3 - Transcrever Arquivo Específico")
		fmt.Println("0 - Sair")
		fmt.Println("------------------------------")

		var opcao string
		fmt.Print("Escolha uma opção: ")
		fmt.Scanln(&opcao)

		switch opcao {
		case "1":
			monitor, err := backend.DetectarMonitor()
			if err != nil {
				fmt.Printf("❌ Erro ao detectar monitor: %v\n", err)
				continue
			}

			// Inicia a gravação interativa (ENTER para parar)
			err = backend.GravarAteParar(monitor)
			if err != nil {
				fmt.Printf("❌ Erro na gravação: %v\n", err)
			}

		case "2":
			// Garante que o Whisper está pronto antes de transcrever
			if err := backend.InstalarWhisper(); err != nil {
				fmt.Printf("❌ Erro com Whisper: %v\n", err)
				continue
			}

			err := backend.TranscreverUltimo()
			if err != nil {
				fmt.Printf("❌ Erro na transcrição: %v\n", err)
			}

		case "3":
			backend.TranscreverArquivo()

		case "0":
			fmt.Println("Encerrando... Até logo!")
			return

		default:
			fmt.Println("⚠️ Opção inválida, tente novamente.")
		}
	}
}
