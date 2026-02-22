//go:build linux

package backend

import (
	"fmt"
	"os/exec"
)

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
