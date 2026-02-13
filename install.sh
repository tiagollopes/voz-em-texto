#!/bin/bash

echo "======================================"
echo " Instalando dependências do projeto "
echo " Voz em Texto (Go + Whisper + GUI) "
echo "======================================"

echo ""
echo "⚙ Atualizando repositórios..."
sudo apt update

echo ""
echo "📦 Instalando dependências principais..."
sudo apt install -y cmake build-essential ffmpeg git pkg-config

echo ""
echo "🖥️ Instalando dependências gráficas (GUI Fyne)..."
sudo apt install -y libgl1-mesa-dev xorg-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev

echo ""
echo "⚙ Verificando Go instalado..."
if ! command -v go >/dev/null 2>&1; then
    echo "❌ Go não encontrado."
    echo "Instale em: https://go.dev/dl/"
    exit 1
fi

echo "✅ Go encontrado."

echo ""
echo "📦 Baixando dependências do projeto..."
go mod tidy

echo ""
echo "======================================"
echo "✅ Instalação concluída!"
echo "======================================"
echo ""
echo "Para rodar CLI:"
echo "   go run main.go"
echo ""
echo "Para rodar GUI:"
echo "   go run gui.go"
echo ""
echo "Ou compilar:"
echo "   go build"
echo ""
