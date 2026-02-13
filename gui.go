package main

import (
	"fmt"
	"time"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"image/color"
	"fyne.io/fyne/v2/canvas"

)

// =========================
// VARIÁVEIS GLOBAIS
// =========================
var gravando bool
var transcrevendo bool

var tempoInicio time.Time

var btnGravar *widget.Button
var btnParar *widget.Button
var btnTranscrever *widget.Button

func popupSelecionarAudio(w fyne.Window, status *widget.Label) {

	lista := listarAudiosInput()

	if len(lista) == 0 {
		dialog.ShowInformation(
			"Input vazio",
			"Nenhum áudio encontrado em /input",
			w,
		)
		return
	}

	selectWidget := widget.NewSelect(lista, nil)

	dialog.ShowCustomConfirm(
		"Selecionar áudio",
		"Transcrever",
		"Cancelar",
		selectWidget,
		func(confirm bool) {

			if confirm && selectWidget.Selected != "" {

				nome := selectWidget.Selected

				status.SetText("⚙️ Transcrevendo: " + nome)

				go func() {

					transcrevendo = true
					atualizarBotoes()

					caminho := "input/" + nome

					duracao, _ := duracaoArquivo(caminho)
					go progressoTranscricao(status, duracao)

					cmd := exec.Command(
						"./whisper/build/bin/whisper-cli",
						"-m", "whisper/models/ggml-base.bin",
						"-f", caminho,
						"-l", "pt",
						"-otxt",
					)

					err := cmd.Run()
					if err != nil {
						status.SetText("❌ Erro na transcrição")
						transcrevendo = false
						atualizarBotoes()
						return
					}

					os.Rename(
						caminho+".txt",
						"output/"+nome+".txt",
					)

					status.SetText("✅ Transcrição concluída")

					transcrevendo = false
					atualizarBotoes()
				}()
			}
		},
		w,
	)
}


func atualizarBotoes() {

	if gravando {

		btnGravar.Disable()
		btnTranscrever.Disable()
		btnParar.Enable()
		return
	}

	if transcrevendo {

		btnGravar.Disable()
		btnParar.Disable()
		btnTranscrever.Disable()
		return
	}

	// Estado parado
	btnGravar.Enable()
	btnTranscrever.Enable()
	btnParar.Disable()
}


// =========================
// ATUALIZAR REC TEMPO REAL
// =========================
func atualizarREC(dot *canvas.Text, label *widget.Label, bar *widget.ProgressBar) {

	for  {

		if !gravando {
		return
		}
		decorrido := time.Since(tempoInicio)

		min := int(decorrido.Minutes())
		seg := int(decorrido.Seconds()) % 60

		// Atualiza tempo (texto nunca muda de posição)
		label.SetText(
			fmt.Sprintf("REC %02d:%02d", min, seg),
		)

		// Piscar mudando intensidade do vermelho
		if seg%2 == 0 {

			// Vermelho forte
			dot.Color = color.RGBA{255, 0, 0, 255}

		} else {

			// Vermelho fraco (quase apagado)
			//dot.Color = color.RGBA{120, 0, 0, 255}
			dot.Color = color.RGBA{255, 0, 0, 80}
		}

		dot.Refresh()

		// Barra animada
		bar.SetValue(float64(seg%10) / 10.0)

		time.Sleep(1 * time.Second)
	}
}

func progressoTranscricao(status *widget.Label, duracao float64) {

	inicio := time.Now()
	tempoEstimado := duracao * 3.0

	mostrouFinalizando := false

	for {

		if !transcrevendo {

			// Mostra finalizando uma vez
			if !mostrouFinalizando {
				status.SetText("⚙️ Finalizando transcrição...")
				time.Sleep(1 * time.Second)
			}

			return
		}

		decorrido := time.Since(inicio).Seconds()
		percent := (decorrido / tempoEstimado) * 100

		if percent >= 95 {
			status.SetText("⚙️ Finalizando transcrição...")
			mostrouFinalizando = true
			time.Sleep(1 * time.Second)
			continue
		}

		texto := fmt.Sprintf(
			"⚙️ Transcrevendo... %.0f%%",
			percent,
		)

		status.SetText(texto)

		time.Sleep(1 * time.Second)
	}
}


func listarAudiosInput() []string {

	files, err := os.ReadDir("input")
	if err != nil {
		return []string{}
	}

	var lista []string

	for _, f := range files {
		lista = append(lista, f.Name())
	}

	return lista
}

// =========================
// MAIN GUI
// =========================
func main() {

	a := app.New()
	w := a.NewWindow("Voz em Texto")
	os.MkdirAll("output", 0755)
	// =========================
	// WIDGETS REC
	// =========================
	// Bolinha
	recDot := canvas.NewText("●", color.RGBA{255, 0, 0, 255})
	recDot.TextSize = 16
	recDot.Hide()

	// Texto REC
	recLabel := widget.NewLabel("REC 00:00")
	recLabel.Hide()

	// Barra
	recBar := widget.NewProgressBar()
	recBar.Hide()

	// =========================
	// STATUS BAR
	// =========================
	status := widget.NewLabel("⚙️")

	// =========================
	// BOTÃO GRAVAR
	// =========================
	btnGravar = widget.NewButton("⚙️ Gravar", func() {


		if gravando {
			status.SetText("⚠️ Já está gravando")
			return
		}

		status.SetText("🔍 Detectando monitor...")

		go func() {

			monitor, err := detectarMonitor()
			if err != nil {
				status.SetText("❌ Erro ao detectar monitor")
				return
			}

			err = iniciarGravacao(monitor)
			if err != nil {
				status.SetText("❌ Erro ao iniciar gravação")
				return
			}

			// Inicia REC
			gravando = true
			atualizarBotoes()

			recDot.Show()
			recLabel.Show()
			recBar.Show()

			tempoInicio = time.Now()

			go atualizarREC(recDot, recLabel, recBar)

			status.SetText("⚙ Gravando...")
		}()
	})

	// =========================
	// BOTÃO PARAR
	// =========================
	btnParar = widget.NewButton("⏹ Parar", func() {

		if !gravando {
			status.SetText("⚠️ Nada gravando")
			return
		}

		status.SetText("⏹ Finalizando gravação...")

		go func() {

			err := pararGravacao()
			if err != nil {
				status.SetText("❌ Erro ao parar")
				return
			}

			gravando = false
			atualizarBotoes()

			recDot.Hide()
			recLabel.Hide()
			recBar.Hide()

			recLabel.SetText("REC 00:00")
			recBar.SetValue(0)

			status.SetText("✅ Áudio salvo")

			dialog.NewConfirm(
				"Transcrever",
				"Deseja transcrever agora?",
				func(resposta bool) {

					if resposta {

						transcrevendo = true
						atualizarBotoes()

						status.SetText("⚙️ Transcrevendo áudio...")

						go func() {

							err := instalarWhisper()
							if err != nil {
								status.SetText("❌ Erro Whisper")
								transcrevendo = false
								atualizarBotoes()
								return
							}

							audioPath := "output/" + ultimoAudioGerado + ".mp3"

							// duração
							duracao, _ := duracaoArquivo(audioPath)

							// inicia progresso
							go progressoTranscricao(status, duracao)

							// transcreve
							err = transcreverUltimo()
							if err != nil {
								status.SetText("❌ Erro transcrição")
								transcrevendo = false
								atualizarBotoes()
								return
							}

							// finaliza
							status.SetText("✅ Transcrição concluída")

							transcrevendo = false
							atualizarBotoes()
						}()
					}
				},
				w, // janela pai
			).Show()
		}()
	})

	// =========================
	// BOTÃO TRANSCRIBIR
	// =========================
	btnTranscrever = widget.NewButton("🧠 Transcrever (input)", func() {
	popupSelecionarAudio(w, status)
	})

	btnSair := widget.NewButton("❌ Sair", func() {
	w.Close()
	})
	btnSair.Importance = widget.LowImportance

	// =========================
	// LAYOUT
	// =========================
	conteudo := container.NewVBox(

		widget.NewLabel("Sistema Voz em Texto"),

		container.NewHBox(
			recDot,
			recLabel,
		),
		recBar,

		btnGravar,
		btnParar,
		btnTranscrever,
		btnSair,
	)

	w.SetContent(
		container.NewBorder(
			nil,
			status, // barra inferior
			nil,
			nil,
			conteudo,
		),
	)

	w.Resize(fyne.NewSize(420, 260))
	w.CenterOnScreen()
	atualizarBotoes()


	w.ShowAndRun()
}
