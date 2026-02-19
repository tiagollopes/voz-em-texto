# 🎙️ Voz em Texto — Go (Linux)

Projeto experimental em Golang para gravação de áudio do sistema e transcrição automática offline utilizando Whisper.cpp.

Desenvolvido e testado em ambiente Linux (Ubuntu/Lubuntu).

#  Arquitetura Modular

O projeto foi refatorado para uma arquitetura em domínios independentes, seguindo boas práticas de organização em Go.

- cmd/ → EntryPoints
- internal/ → Domínios de negócio

##  Funcionalidades

- **Gravação de áudio do sistema**: Captura o áudio interno via PulseAudio monitor utilizando FFmpeg.
- **Transcrição Offline**: Integração com Whisper.cpp para processamento local.
- **Portabilidade**: Uso de `bundled.go` para embutir ícones, evitando caminhos quebrados ao mover o executável.
- **Organização de Arquivos**:

    - `audio/`: Arquivos temporários.
    - `input/`: Para áudios externos.
    - `output/`: Resultados finais em MP3 e TXT.

## 🛠️ Instalação e Dependências

**1. Dependências do Sistema**

O script `install.sh` automatiza a instalação de:

- `cmake`, `ffmpeg`, `build-essential`, `pkg-config`.
- Dependências X11 para a interface gráfica Fyne.

## Estrutura

- **cmd/gui/** → Interface gráfica (Fyne)
- **cmd/cli/** → Interface terminal
- **internal/audio/** → Captura e gravação de áudio
- **internal/transcribe/** → Execução Whisper e IA
- **internal/progress/** → Feedback visual de progresso
- **internal/system/** → Infraestrutura e paths
- **internal/backend/** → Orquestração leve e dependências

# Funcionalidades

### 🎧 Gravação de áudio do sistema

Captura áudio interno via PulseAudio monitor usando FFmpeg.

### Transcrição Offline

Processamento local com Whisper.cpp (sem nuvem).

### Arquitetura desacoplada

IA, captura e feedback separados por domínio.

### 🖥️ Interfaces disponíveis

- GUI (Fyne)
- CLI (Terminal)

###  Organização de Arquivos

- `audio/` → Temporários de gravação
- `input/` → Áudios externos
- `output/` → Resultados finais (.mp3 / .txt)

# 🛠️ Instalação

## Dependências

Script automático:

<pre>chmod +x install.sh
./install.sh</pre>

Instala:

- cmake
- ffmpeg
- build-essential
- pkg-config
- dependências gráficas Fyne

# ▶️ Como Executar

## GUI

<pre>go run ./cmd/gui</pre>

## CLI

<pre>go run ./cmd/cli</pre>

# Fluxos de Trabalho

### Gravar + Transcrever

Grava o áudio e inicia a transcrição automaticamente.

### Transcrever Externo

Seleciona arquivo da pasta `input/` e gera `.txt` em `output/`.

# Status do Projeto

- Arquitetura modular concluída
- IA isolada no domínio transcribe
- Backend limpo
- GUI e CLI desacoplados
- Execução 100% offline

# Licença

Este projeto é de uso livre para:

- Estudos

- Modificações

- Uso pessoal ou comercial

Peço apenas que mantenha os créditos ao autor original:

Tiago Lopes

GitHub: https://github.com/tiagollopes


***Projeto experimental em Golang para automação de voz → texto offline.***

# Autor

Feito por **Tiago LLopes** - Santos/SP - Brasil  🇧🇷
