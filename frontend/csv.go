package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"magic-tournament/logic"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: filetool input.txt output.txt")
		return
	}

	infile := os.Args[1]
	outfile := os.Args[2]

	in, err := os.Open(infile)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	t := logic.NewTournament()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			continue
		}

		t.AddMatch(logic.MatchResult{
			Player1:   parts[0],
			Player2:   parts[1],
			GamesWon1: toInt(parts[2]),
			GamesWon2: toInt(parts[3]),
		})
	}

	out, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	out.WriteString(logic.FormatStandings(t))
}

func toInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
