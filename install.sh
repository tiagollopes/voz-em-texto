#!/bin/bash

echo "======================================"
echo " Instalando dependências do projeto "
echo "======================================"

echo ""
echo "Atualizando repositórios..."
sudo apt update

echo ""
echo "Instalando pacotes necessários..."
sudo apt install -y \
    cmake \
    build-essential \
    ffmpeg \
    git \
    pkg-config

echo ""
echo "✅ Dependências instaladas!"
echo ""
echo "Agora rode:"
echo "go run main.go"
