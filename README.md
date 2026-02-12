#  🎙️ Voz em Texto — CLI em Go (Linux)

Projeto experimental em **Golang** para gravação de áudio do sistema e transcrição automática em texto usando **Whisper.cpp**.

Desenvolvido e testado em ambiente Linux (Lubuntu/Ubuntu).

##  Objetivo

Este é um projeto de estudos/testes para:

- Captura de áudio do PC

- Transcrição offline

- Automação via CLI

- Integração Go + FFmpeg + Whisper.cpp

Não é um produto final — está em evolução contínua.

##  ⚙️ Requisitos

Antes de rodar, o sistema precisa ter:

- Linux (Ubuntu / Lubuntu recomendado)

- Go instalado

- Permissão sudo

As demais dependências são instaladas automaticamente.

##  📦 Instalação

Clone o repositório:

<pre>git clone https://github.com/tiagollopes/voz-em-texto.git</pre>
<pre>cd voz-em-texto</pre>

Dê permissão ao instalador:

<pre>chmod +x install.sh</pre>

Execute:

<pre>./install.sh</pre>

O script instala:

- ffmpeg

- cmake

- build-essential

- git

- whisper.cpp

- modelo de transcrição

##  ▶️ Execução

Rodar o sistema:

<pre>go run main.go</pre>

##   Menu do sistema

1 - Gravar áudio

2 - Transcrever áudio existente

0 - Sair

##  Estrutura de pastas

voz-em-texto/

├── main.go

├── install.sh

├── audio/     → gravação temporária

├── input/     → áudios para transcrever

├── output/    → resultados finais

└── whisper/   → instalado automaticamente

##  Como funciona

**Gravação**

- Captura áudio do monitor do sistema

- Salva com timestamp

- Copia para /output

***Exemplo:***

<pre>gravacao_12022026_173404.mp3</pre>

**Transcrição**

- Usa Whisper.cpp offline

- Idioma: Português

- Gera .txt com mesmo nome

***Exemplo:***

<pre>gravacao_12022026_173404.txt</pre>

**Licença**

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

**Feito por Tiago LLopes** - Santos/SP - Brasil  🇧🇷
