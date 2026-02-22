package main

import (
	"voz-em-texto/internal/audio"
	"voz-em-texto/internal/transcribe"
	"fmt"
	"time"
	"os"
	"os/exec"
	"path/filepath"
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
var cmdGravacao *exec.Cmd

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

					// progresso GUI
					duracao, _ := transcribe.DuracaoArquivo(caminho)
					go progressoTranscricao(status, duracao)

					// chama domínio IA
					err := transcribe.TranscreverCaminho(caminho)
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

					// sucesso
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
	// Mantendo sua lógica de 3x a duração para ser bem brando
	tempoEstimado := duracao * 3.0
	// Se for um áudio minúsculo, garante ao menos 3s de barra
	if tempoEstimado < 3 { tempoEstimado = 3 }

	mostrouFinalizando := false

	for {
		// 1. Checa primeiro se o motor parou antes de atualizar a tela
		if !transcrevendo {
			if cancelado {
				fyne.Do(func() { status.SetText("⛔ Transcrição cancelada") })
			} else {
				fyne.Do(func() { status.SetText("✅ Transcrição concluída") })
			}
			return
		}

		decorrido := time.Since(inicio).Seconds()
		percent := (decorrido / tempoEstimado) * 100

		if percent > 95 {
			percent = 95
		}

		if percent >= 95 {
			if !mostrouFinalizando {
				fyne.Do(func() { status.SetText("⚙️ Finalizando transcrição...") })
				mostrouFinalizando = true
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// 2. Aqui está o pulo do gato: usamos 500ms para garantir que o
		// usuário veja pelo menos um número antes de terminar.
		texto := fmt.Sprintf("⚙️ Transcrevendo... %.0f%%", percent)
		fyne.Do(func() {
			status.SetText(texto)
		})

		time.Sleep(500 * time.Millisecond)
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

			monitor, err := audio.DetectarMonitor()
			if err != nil {
				status.SetText("❌ Erro ao detectar monitor")
				return
			}

			err = audio.IniciarGravacao(monitor)
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

			err := audio.PararGravacao()
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

						// 1. Mensagem inicial no rodapé
						fyne.Do(func() {
							status.SetText("⚙️ Transcrevendo áudio...")
						})

						go func() {
							// --- ACRESCIMO DE SEGURANÇA PARA WINDOWS ---
							time.Sleep(500 * time.Millisecond)

							err := transcribe.InstalarWhisper()
							if err != nil {
								// POPUP DE ERRO: Só sai da tela quando você der OK
								fyne.Do(func() {
									dialog.ShowError(fmt.Errorf("Whisper não localizado em bin/windows"), w)
									status.SetText("❌ Erro Whisper")
								})
								transcrevendo = false
								cancelado = true
								atualizarBotoes()
								return
							}

							audioPath := audio.UltimoAudioGerado
							if audioPath == "" {
								audioPath = filepath.Join("output", "audio_recente.wav")
							}

							nomeBase := filepath.Base(audioPath)

							fyne.Do(func() {
							status.SetText("⚙️ Transcrevendo: " + nomeBase)
							})

							// Espera um pouco para o usuário ler o popup antes de iniciar a barra
							//time.Sleep(1000 * time.Millisecond)

							// duração
							duracao, _ := transcribe.DuracaoArquivo(audioPath)

							// inicia progresso (Este loop vai atualizar o status com %)
							go progressoTranscricao(status, duracao)

							// transcreve
							err = transcribe.TranscreverUltimo()
							if err != nil {
								transcrevendo = false
								fyne.Do(func() {
									// POPUP DE ERRO REAL: Vai mostrar por que a transcrição falhou
									dialog.ShowError(fmt.Errorf("Erro na transcrição:\n%v", err), w)
									status.SetText(fmt.Sprintf("❌ Erro: %v", err))
								})
								cancelado = true
								atualizarBotoes()
								return
							}

							// finaliza
							transcrevendo = false
							fyne.Do(func() {
								status.SetText("✅ Transcrição concluída: " + nomeBase)
								dialog.ShowInformation("Sucesso", "A transcrição foi finalizada com sucesso!", w)
							})

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

		err := transcribe.PararTranscricao()
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
