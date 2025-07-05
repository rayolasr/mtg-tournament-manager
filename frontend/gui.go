package main

import (
	"fmt"
	"magic-tournament/logic"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	fmt.Println("Iniciando aplicación de torneo de Magic...")

	// Agregar timeout para debug
	go func() {
		time.Sleep(30 * time.Second)
		fmt.Println("La aplicación está iniciando... esto puede tomar varios minutos en la primera ejecución")
	}()

	a := app.New()
	w := a.NewWindow("Torneo de Magic")

	tournament := logic.NewTournament()

	// Entradas de usuario
	player1 := widget.NewEntry()
	player2 := widget.NewEntry()
	gamesWon1 := widget.NewEntry()
	gamesWon2 := widget.NewEntry()

	form := container.NewVBox(
		widget.NewLabel("Introduce una partida:"),
		widget.NewForm(
			widget.NewFormItem("Jugador 1", player1),
			widget.NewFormItem("Jugador 2", player2),
			widget.NewFormItem("Games ganados J1", gamesWon1),
			widget.NewFormItem("Games ganados J2", gamesWon2),
		),
	)

	output := widget.NewMultiLineEntry()
	output.SetMinRowsVisible(15)
	output.SetPlaceHolder("Clasificación y resumen aparecerán aquí")

	addMatch := widget.NewButton("Añadir partida", func() {
		gw1, err1 := parseInt(gamesWon1.Text)
		gw2, err2 := parseInt(gamesWon2.Text)

		if err1 != nil || err2 != nil || player1.Text == "" || player2.Text == "" {
			dialog.ShowError(fmt.Errorf("Datos inválidos"), w)
			return
		}

		match := logic.MatchResult{
			Player1:   player1.Text,
			Player2:   player2.Text,
			GamesWon1: gw1,
			GamesWon2: gw2,
		}
		tournament.AddMatch(match)
		// tournament.UpdateChampion()
		showResults(tournament, output)

		// limpiar entradas
		player1.SetText("")
		player2.SetText("")
		gamesWon1.SetText("")
		gamesWon2.SetText("")
	})

	content := container.NewVBox(
		form,
		addMatch,
		widget.NewSeparator(),
		widget.NewLabel("Resultados:"),
		output,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(700, 600))
	w.ShowAndRun()
	fmt.Println("Configuración completada, mostrando ventana...")
	w.ShowAndRun() // Esta línea bloquea hasta que se cierre la aplicación
	fmt.Println("Aplicación cerrada correctamente")
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
	return n, err
}

func showResults(t *logic.Tournament, out *widget.Entry) {
	text := logic.FormatStandings(t)
	out.SetText(text)
}
