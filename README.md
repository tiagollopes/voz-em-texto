#  🎙️ Voz em Texto — Go (Linux)

Projeto experimental em **Golang** para gravação de áudio do sistema e transcrição automática em texto usando **Whisper.cpp**.

Desenvolvido e testado em ambiente Linux (Lubuntu/Ubuntu).

Possui:

- CLI (terminal)

- Interface gráfica (Fyne)

- Gravação em tempo real

- Transcrição automática ou manual

- Suporte a arquivos externos

# Funcionalidades

**Gravação de áudio do sistema**

- Captura áudio interno via PulseAudio monitor

- Usa ffmpeg

- Grava até clicar em parar

- Salva automaticamente em:

<pre>output/gravacao_ddmmaaaa_hhmmss.mp3</pre>

# Transcrição com Whisper.cpp

- Instala Whisper automaticamente

- Compila via make

- Baixa modelo (tiny / base / small…)

Força idioma português:

<pre>-l pt</pre>

Saída:

<pre>output/gravacao_xxx.txt</pre>
# Transcrição de arquivos externos

Coloque áudios em:

<pre>input/</pre>

A GUI lista os arquivos disponíveis.

Fluxo:

Seleciona → Transcreve → Salva em:

<pre>output/nome.txt</pre>

# 🖥️ Interface Gráfica

Desenvolvida com **Fyne v2**.

**Elementos**

- Botão Gravar

- Botão Parar

- Botão Transcrever (input)

- Botão Sair

**Indicador REC**

Durante gravação:

- Bolinha vermelha piscando

- Tempo decorrido

- Barra animada

Exemplo visual:

● REC 00:32
██████░░░░

Características:

- Só aparece gravando

- Some ao parar

- Não move layout

# 📊 Progresso de Transcrição

- Barra baseada na duração do áudio

- Vai até 95%

- Depois:

Finalizando transcrição…

# Fluxos do Sistema

-  1️⃣  Gravar + transcrever

Gravar → Parar → Popup → Sim → Transcreve

-  2️⃣  Gravar sem transcrever

Gravar → Parar → Popup → Não

Áudio fica salvo em output/.

-  3️⃣  Transcrever arquivos externos

Botão → Lista → Seleciona → Transcreve

# Estrutura do Projeto

voz-em-texto/

│

├── main.go        → CLI

├── gui.go         → Interface gráfica

├── backend.go     → Funções de gravação/transcrição

├── install.sh     → Instalador Linux

│

├── audio/         → Áudio temporário

├── input/         → Áudios externos

├── output/        → Resultados

│

└── whisper/       → Whisper.cpp (auto instalado)

# Dependências Linux

Instaladas via install.sh:

- cmake

- build-essential

- ffmpeg

- git

- pkg-config

- libgl1-mesa-dev

- xorg-dev

- libxcursor-dev

- libxrandr-dev

- libxinerama-dev

- libxi-dev

# Dependências Go

<pre>fyne.io/fyne/v2</pre>

Instalar:

<pre>go mod tidy</pre>

# ⚙️ Instalação

-  1️⃣  Clonar

<pre>git clone https://github.com/tiagollopes/voz-em-texto.git </pre>

<pre>cd voz-em-texto</pre>

-  2️⃣  Rodar instalador

<pre>chmod +x install.sh</pre>

<pre>./install.sh</pre>

O script:

- Instala dependências

- Clona Whisper.cpp

- Compila

- Baixa modelo

# ▶️ Executar

**GUI**

<pre>go run gui.go backend.go</pre>

ou build:

<pre>go build -o voz-em-texto</pre>

<pre>./voz-em-texto</pre>

**CLI**

<pre>go run main.go</pre>

# Estado do Projeto

- Área	Status

- CLI	✅

- GUI	✅

- Gravação	✅

- Transcrição	✅

- Indicador REC	✅

- Progresso	✅

# Licença

Este projeto é de uso livre para:

- Estudos

- Modificações

-- Uso pessoal ou comercial

Peço apenas que mantenha os créditos ao autor original:

Tiago Lopes

GitHub: https://github.com/tiagollopes

##  Status do projeto

- Em desenvolvimento / testes

Funcionalidades podem mudar ou evoluir.

##  Futuras melhorias

- Resumo automático de texto

- Tradução de transcrição

- Interface gráfica

- Exportação em PDF

- Batch de arquivos

##  Contribuição

Sinta-se livre para:

- Abrir issues

- Sugerir melhorias

- Fazer fork do projeto

***Projeto experimental em Golang para automação de voz → texto offline.***

# Autor

Feito por **Tiago LLopes** - Santos/SP - Brasil  🇧🇷
