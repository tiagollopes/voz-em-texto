package backend

import (
	"fmt"
	"os/exec"
)

// ===============================
// Checagem Dependencia Linux
// ===============================
func ChecarDependencias() error {

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




