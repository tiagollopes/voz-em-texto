#!/bin/bash

echo "======================================"
echo " Instalador - Voz em Texto "
echo "======================================"
echo ""

echo "⚙ Preparando ambiente..."

sudo apt update > /dev/null 2>&1

echo "📦 Instalando dependências principais..."
sudo apt install -y cmake build-essential ffmpeg git pkg-config pulseaudio-utils > /dev/null 2>&1

echo "🖥️ Instalando dependências gráficas..."
sudo apt install -y libgl1-mesa-dev xorg-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev > /dev/null 2>&1

echo "⚙ Ajustando permissões dos executáveis..."

find . -maxdepth 1 -type f -name "voz*" -exec chmod +x {} ; 2>/dev/null

echo ""
echo "======================================"
echo "✅ Instalação concluída com sucesso!"
echo "======================================"
echo ""
echo "Execute o aplicativo normalmente:"
echo "   ./vozgui"
echo ""
