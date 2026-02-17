package main

import (
	"voz-em-texto/internal/backend"

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
	"fyne.io/fyne/v2/theme"

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
var btnPararTranscricao *widget.Button
var cancelado bool


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
				fyne.Do(func() {
					status.SetText("⚙️ Transcrevendo: " + nome)
				})
				go func() {

					transcrevendo = true
					cancelado = false
					fyne.Do(func() {
					atualizarBotoes()
					})

					caminho := "input/" + nome

					duracao, _ := backend.DuracaoArquivo(caminho)
					go progressoTranscricao(status, duracao)

					cmd := exec.Command(
						"./whisper/build/bin/whisper-cli",
						//"-m", "whisper/models/ggml-base.bin",
						"-m", "whisper/models/ggml-tiny.bin",
						"-f", caminho,
						"-l", "pt",
						"-otxt",
					)

					err := cmd.Run()
					if err != nil {
						fyne.Do(func() {
						status.SetText("❌ Erro na transcrição")
						})
						transcrevendo = false
						cancelado = true
						fyne.Do(func() {
						atualizarBotoes()
						})
						return
					}

					os.Rename(
						caminho+".txt",
						"output/"+nome+".txt",
					)
					fyne.Do(func() {
					status.SetText("✅ Transcrição concluída")
					})
					transcrevendo = false
					cancelado = false
					fyne.Do(func() {
					atualizarBotoes()
					})
				}()
			}
		},
		w,
	)
}


func atualizarBotoes() {

	fyne.Do(func() {

		if gravando {

			btnGravar.Disable()
			btnTranscrever.Disable()
			btnParar.Enable()
			btnPararTranscricao.Hide()
			return
		}

		if transcrevendo {

			btnGravar.Disable()
			btnParar.Disable()
			btnTranscrever.Disable()
			btnPararTranscricao.Show()
			btnPararTranscricao.Enable()
			return
		}

		// parado
		btnGravar.Enable()
		btnTranscrever.Enable()
		btnParar.Disable()
		btnPararTranscricao.Hide()
	})
}



// =========================
// ATUALIZAR REC TEMPO REAL
// =========================
func atualizarREC(dot *canvas.Text, label *widget.Label, bar *widget.ProgressBar) {

	for {

		if !gravando {
			return
		}

		decorrido := time.Since(tempoInicio)

		min := int(decorrido.Minutes())
		seg := int(decorrido.Seconds()) % 60

		fyne.Do(func() {

			label.SetText(
				fmt.Sprintf("REC %02d:%02d", min, seg),
			)

			if seg%2 == 0 {
				dot.Color = color.RGBA{255, 0, 0, 255}
			} else {
				dot.Color = color.RGBA{255, 0, 0, 80}
			}

			dot.Refresh()
			bar.SetValue(float64(seg%10) / 10.0)
		})

		time.Sleep(1 * time.Second)
	}
}

func progressoTranscricao(status *widget.Label, duracao float64) {

	inicio := time.Now()
	tempoEstimado := duracao * 3.0

	mostrouFinalizando := false

	for {

		// Sai se terminou ou cancelou
		if !transcrevendo {

			if cancelado {

				fyne.Do(func() {
					status.SetText("⛔ Transcrição cancelada")
				})

			} else {

				fyne.Do(func() {
					status.SetText("✅ Transcrição concluída")
				})
			}

			return
		}

		decorrido := time.Since(inicio).Seconds()
		percent := (decorrido / tempoEstimado) * 100

		// LIMITADOR — mantém 95%
		if percent > 95 {
			percent = 95
		}

		// Quando chega em 95 → entra em modo finalização
		if percent >= 95 {

			if !mostrouFinalizando {

				fyne.Do(func() {
					status.SetText("⚙️ Finalizando transcrição...")
				})

				mostrouFinalizando = true
			}

			// 🔒 trava aqui até terminar ou cancelar
			time.Sleep(500 * time.Millisecond)
			continue
		}

		texto := fmt.Sprintf(
			"⚙️ Transcrevendo... %.0f%%",
			percent,
		)

		fyne.Do(func() {
			status.SetText(texto)
		})

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

	a.SetIcon(theme.MediaPlayIcon())

	a.Settings().SetTheme(theme.LightTheme())

	w := a.NewWindow("Voz em Texto")

	w.SetIcon(theme.MediaPlayIcon())


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

			monitor, err := backend.DetectarMonitor()
			if err != nil {
				status.SetText("❌ Erro ao detectar monitor")
				return
			}

			err = backend.IniciarGravacao(monitor)
			if err != nil {
				status.SetText("❌ Erro ao iniciar gravação")
				return
			}

			// Inicia REC
			gravando = true
			atualizarBotoes()
			fyne.Do(func() {
				recDot.Show()
				recLabel.Show()
				recBar.Show()
			})

			tempoInicio = time.Now()

			go atualizarREC(recDot, recLabel, recBar)
			fyne.Do(func() {
			status.SetText("⚙ Gravando...")
			})
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

			err := backend.PararGravacao()
			if err != nil {
				status.SetText("❌ Erro ao parar")
				return
			}

			gravando = false
			atualizarBotoes()
			fyne.Do(func() {
				recDot.Hide()
				recLabel.Hide()
				recBar.Hide()

				recLabel.SetText("REC 00:00")
				recBar.SetValue(0)

				status.SetText("✅ Áudio salvo")
			})
			dialog.NewConfirm(
				"Transcrever",
				"Deseja transcrever agora?",
				func(resposta bool) {

					if resposta {

						transcrevendo = true
						cancelado = false
						atualizarBotoes()
						fyne.Do(func() {
						status.SetText("⚙️ Transcrevendo áudio...")
						})
						go func() {

							err := backend.InstalarWhisper()
							if err != nil {
								fyne.Do(func() {
								status.SetText("❌ Erro Whisper")
								})
								transcrevendo = false
								cancelado = true
								atualizarBotoes()
								return
							}

							audioPath := "output/" + backend.UltimoAudioGerado + ".mp3"

							// duração
							duracao, _ := backend.DuracaoArquivo(audioPath)

							// inicia progresso
							go progressoTranscricao(status, duracao)

							// transcreve
							err = backend.TranscreverUltimo()
							if err != nil {
								fyne.Do(func() {
								status.SetText("❌ Erro transcrição")
								})
								transcrevendo = false
								cancelado = true
								atualizarBotoes()
								return
							}

							// finaliza
							fyne.Do(func() {
							status.SetText("✅ Transcrição concluída")
							})

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
	btnTranscrever = widget.NewButton(" Transcrever (input)", func() {
	popupSelecionarAudio(w, status)
	})

	btnSair := widget.NewButton("❌ Sair", func() {
	w.Close()
	})
	btnSair.Importance = widget.LowImportance

	// =========================
	// BOTÃO PARAR TRANSCRIÇÃO
	// =========================
	btnPararTranscricao = widget.NewButton("⛔ Parar Transcrição", func() {

		err := backend.PararTranscricao()
		if err != nil {
			//status.SetText("❌ Erro ao parar transcrição")
			status.SetText(err.Error())
			return
		}

		status.SetText("⛔ Transcrição cancelada")

		transcrevendo = false
		cancelado = true
		atualizarBotoes()
	})
	btnPararTranscricao.Hide()

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
		btnPararTranscricao,
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
