# 🎙️ Voz em Texto — Go (Multiplataforma)

Projeto experimental em Golang para gravação de áudio do sistema e transcrição automática offline utilizando **Whisper.cpp**. O sistema agora é totalmente multiplataforma, suportando **Linux** e **Windows** tanto em interface de linha de comando (CLI) quanto gráfica (GUI).

## Novidades: Suporte Windows

O projeto foi atualizado para rodar nativamente em Windows. Para garantir o funcionamento, é necessário utilizar a estrutura da pasta <pre>`bin/`</pre> para dependências externas.

---

## 🖥️ Interfaces Disponíveis

* **GUI (Fyne):** Interface gráfica amigável para gravação e transcrição.
* **CLI (Terminal):** Versão leve para uso via linha de comando.

---

## 🛠️ Instalação e Dependências

### Linux

O projeto foi desenvolvido e testado em ambiente Linux (Ubuntu/Lubuntu).

1. **Dependências do Sistema:**

   O script <pre>`install.sh`</pre> automatiza a instalação de: <pre>`cmake`</pre>, <pre>`ffmpeg`</pre>, <pre>`build-essential`</pre>, <pre>`pkg-config`</pre> e dependências X11 para a interface gráfica Fyne.

   <pre>chmod +x install.sh
   ./install.sh</pre>

### 🪟 Windows

Para rodar no Windows, o sistema depende de binários e bibliotecas específicas localizadas na pasta bin/.

Dependências Obrigatórias:

Devido ao tamanho, alguns arquivos devem ser baixados na aba Releases deste repositório:

Coloque <pre>ffmpeg.exe</pre> e <pre>ffprobe.exe</pre> em: bin/windows/

Coloque o modelo <pre>ggml-tiny.bin</pre> em: <pre>bin/models/</pre>

As DLLs essenciais (SDL2.dll, whisper.dll, etc.) já estão incluídas no repositório na pasta <pre>bin/windows/</pre>.

## 🏗️ Compilação (Build)

Se você deseja gerar os executáveis manualmente, utilize os comandos abaixo:

###Para Windows (Cross-compilation no Linux)

GUI (Sem janela de terminal):

<pre>
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H=windowsgui -s -w" -o voz-gui.exe ./cmd/gui
</pre>

CLI:

<pre>
GOOS=windows GOARCH=amd64 go build -o voz-cli.exe ./cmd/cli
</pre>

###Para Linux

GUI:

<pre>
go build -o teste-gui-linux ./cmd/gui
</pre>

CLI:

<pre>
go build -o voz-cli-linux ./cmd/cli
</pre>

## Organização de Arquivos e Domínios

O projeto segue uma arquitetura modular baseada em domínios independentes:

- cmd/ → EntryPoints (GUI e CLI).

- internal/audio/ → Captura e gravação de áudio (PulseAudio/Linux e WASAPI/Windows).

- internal/transcribe/ → Execução do Whisper e motor de IA.

- internal/system/ → Gestão de caminhos (Paths) e infraestrutura.

- audio/ → Arquivos temporários de gravação.

- input/ → Para áudios externos que deseja transcrever.

- output/ → Resultados finais em .mp3 e .txt.

## ✨ Funcionalidades


* **🎧 Gravação do Sistema:** Captura áudio interno (o que você ouve) sem necessidade de microfone externo.


* **Transcrição 100% Offline:** Processamento local com Whisper.cpp.


* **📦 Portabilidade:** Uso de bundled.go para embutir ativos (ícones), evitando caminhos quebrados.

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
