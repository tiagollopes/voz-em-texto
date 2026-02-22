package backend

import "fmt"

func Menu() string {

	var opcao string

	fmt.Println("")
	fmt.Println("1 - Gravar áudio")
	fmt.Println("2 - Transcrever áudio existente")
	fmt.Println("0 - Sair")
	fmt.Print("Escolha: ")

	fmt.Scanln(&opcao)

	return opcao
}
