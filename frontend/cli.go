package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"magic-tournament/logic"
)

func main() {
	t := logic.NewTournament()

	fmt.Println("Introduce resultados (formato: jugador1,jugador2,g1,g2). Ctrl+D para terminar:")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			fmt.Println("Formato inválido. Usa jugador1,jugador2,2,1")
			continue
		}

		t.AddMatch(logic.MatchResult{
			Player1:   parts[0],
			Player2:   parts[1],
			GamesWon1: toInt(parts[2]),
			GamesWon2: toInt(parts[3]),
		})
	}

	fmt.Println("\nClasificación:")
	fmt.Print(logic.FormatStandings(t))
}

func toInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
