# 🎙️ Voz em Texto — Go (Linux)

Projeto experimental em **Golang** para gravação de áudio do sistema e transcrição automática em texto usando **Whisper.cpp**.

Desenvolvido e testado em ambiente Linux (Ubuntu/Lubuntu).

## Estrutura Modular
O projeto foi reestruturado para seguir as melhores práticas de organização em Go:
- **`cmd/gui/`**: Ponto de entrada da Interface Gráfica (Fyne).
- **`cmd/cli/`**: Ponto de entrada da Interface de Terminal.
- **`internal/backend/`**: Lógica centralizada para controle de áudio, FFmpeg e Whisper.
- **`assets/`**: Repositório de recursos visuais e ícones.
- **`bundled.go`**: Recursos embutidos (ícones) para garantir portabilidade total do binário.

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

**2. Configuração**

<pre>
git clone [https://github.com/tiagollopes/voz-em-texto.git](https://github.com/tiagollopes/voz-em-texto.git)
cd voz-em-texto
chmod +x install.sh
./install.sh
</pre>

### Como Executar

**Interface Gráfica (GUI)**

<pre>
go run ./cmd/gui
</pre>

**Terminal (CLI)**

<pre>
go run ./cmd/cli
</pre>

**Compilar Executável Único**

<pre>
go build -o voz-em-texto ./cmd/gui
</pre>

### 📊 Fluxos de Trabalho

- Gravar + Transcrever: Grava o áudio, encerra e inicia automaticamente a transcrição Whisper.

- Transcrever Externo: Seleciona um arquivo da pasta input/ e gera o .txt correspondente na output/.

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
