package logic

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// Representa el resultado de un match entre dos jugadores
type MatchResult struct {
	Player1   string
	Player2   string
	GamesWon1 int
	GamesWon2 int
	Champion  string // Nombre del jugador que tiene el cinturón tras este match
}

// Estadísticas de cada jugador
type PlayerStats struct {
	Name      string
	Matches   int
	Wins      int
	Draws     int
	Losses    int
	GamesWon  int
	GamesLost int
	Points    int
	Opponents map[string]bool // Para calcular OMW%
}

// Estructura principal del torneo
type Tournament struct {
	Players  map[string]*PlayerStats
	Matches  []MatchResult
	Champion string // Nombre del jugador que tiene el cinturón actualmente
}

// Inicializa el sistema de log para escribir en un archivo (opcional)
func init() {
	f, err := os.OpenFile("tournament.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("No se pudo abrir el archivo de log, usando stderr")
		return
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Crea un nuevo torneo vacío
func NewTournament() *Tournament {
	log.Println("Nuevo torneo creado")
	return &Tournament{Players: make(map[string]*PlayerStats)}
}

// Añade un match al torneo y actualiza estadísticas y campeón
func (t *Tournament) AddMatch(result MatchResult) {
	log.Printf("Añadiendo match: %s vs %s (%d-%d)\n", result.Player1, result.Player2, result.GamesWon1, result.GamesWon2)
	// Guardar el match
	t.Matches = append(t.Matches, result)

	p1 := t.getOrCreate(result.Player1)
	p2 := t.getOrCreate(result.Player2)

	// Actualizar estadísticas de partidas jugadas y juegos ganados/perdidos
	p1.Matches++
	p2.Matches++
	p1.GamesWon += result.GamesWon1
	p1.GamesLost += result.GamesWon2
	p2.GamesWon += result.GamesWon2
	p2.GamesLost += result.GamesWon1

	// Registrar oponentes para OMW%
	if p1.Opponents == nil {
		p1.Opponents = make(map[string]bool)
	}
	if p2.Opponents == nil {
		p2.Opponents = make(map[string]bool)
	}
	p1.Opponents[p2.Name] = true
	p2.Opponents[p1.Name] = true

	// Actualizar puntos y resultados
	switch {
	case result.GamesWon1 > result.GamesWon2:
		p1.Wins++
		p2.Losses++
		p1.Points += 3
		log.Printf("%s gana el match y suma 3 puntos", p1.Name)
	case result.GamesWon1 < result.GamesWon2:
		p2.Wins++
		p1.Losses++
		p2.Points += 3
		log.Printf("%s gana el match y suma 3 puntos", p2.Name)
	default:
		p1.Draws++
		p2.Draws++
		p1.Points++
		p2.Points++
		log.Printf("Empate entre %s y %s, ambos suman 1 punto", p1.Name, p2.Name)
	}

	// Actualizar el campeón
	t.UpdateChampion(result)

	// Guardar el campeón en el último match registrado
	if len(t.Matches) > 0 {
		t.Matches[len(t.Matches)-1].Champion = t.Champion
		log.Printf("El campeón tras este match es: %s", t.Champion)
	}
}

// Lógica para actualizar el campeón del torneo según el resultado del match
func (t *Tournament) UpdateChampion(match MatchResult) {
	// Si no hay campeón aún, lo asignamos al ganador del primer match (si no hay empate)
	if t.Champion == "" {
		if match.GamesWon1 > match.GamesWon2 {
			t.Champion = match.Player1
			match.Champion = match.Player1
			log.Printf("Primer campeón: %s", t.Champion)
		} else if match.GamesWon2 > match.GamesWon1 {
			t.Champion = match.Player2
			match.Champion = match.Player2
			log.Printf("Primer campeón: %s", t.Champion)
		}
		return
	}

	// Si el campeón está jugando y perdió, el cinturón cambia
	if match.Player1 == t.Champion {
		if match.GamesWon2 > match.GamesWon1 {
			t.Champion = match.Player2
			match.Champion = match.Player2
			log.Printf("Nuevo campeón: %s", t.Champion)
		}
	} else if match.Player2 == t.Champion {
		if match.GamesWon1 > match.GamesWon2 {
			t.Champion = match.Player1
			match.Champion = match.Player1
			log.Printf("Nuevo campeón: %s", t.Champion)
		}
	}
	// Empates no cambian el cinturón
}

// Devuelve el puntero a las estadísticas del jugador, creándolo si no existe
func (t *Tournament) getOrCreate(name string) *PlayerStats {
	if p, ok := t.Players[name]; ok {
		return p
	}
	log.Printf("Nuevo jugador registrado: %s", name)
	t.Players[name] = &PlayerStats{Name: name, Opponents: make(map[string]bool)}
	return t.Players[name]
}

// Devuelve un slice con las estadísticas de todos los jugadores
func (t *Tournament) Standings() []PlayerStats {
	var result []PlayerStats
	for _, p := range t.Players {
		result = append(result, *p)
	}
	return result
}

// Porcentaje de victorias del jugador
func (p PlayerStats) WinPercentage() float64 {
	if p.Matches == 0 {
		return 0
	}
	return float64(p.Wins) / float64(p.Matches) * 100
}

// Calcula el OMW% (Opponent Match Win Percentage) de un jugador
func (t *Tournament) OMW(player *PlayerStats) float64 {
	var total float64
	var count int

	for opponentName := range player.Opponents {
		opponent := t.Players[opponentName]
		if opponent == nil || opponent.Matches == 0 {
			total += 33.33 // Valor por defecto si el oponente tiene 0 partidas
		} else {
			total += float64(opponent.Wins) / float64(opponent.Matches) * 100
		}
		count++
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// Devuelve una cadena con la clasificación y el resumen de partidas
func FormatStandings(t *Tournament) string {
	var sb strings.Builder

	// Encabezado de clasificación
	sb.WriteString("Clasificación de la Liga de Magic\n")
	sb.WriteString("--------------------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("%-3s %-15s %5s %3s %3s %3s %3s %7s %7s\n",
		"#", "Jugador", "Pts", "PJ", "V", "E", "D", "%Vict", "OMW%"))
	sb.WriteString("--------------------------------------------------------------\n")

	// Obtener standings ordenados por puntos y OMW%
	standings := t.Standings()
	sort.SliceStable(standings, func(i, j int) bool {
		if standings[i].Points == standings[j].Points {
			return t.OMW(&standings[i]) > t.OMW(&standings[j])
		}
		return standings[i].Points > standings[j].Points
	})

	// Imprimir la tabla de clasificación
	for i, p := range standings {
		displayName := p.Name
		if p.Name == t.Champion {
			displayName = displayName + "*"
		} else {
			displayName = displayName + " " // Espacio para alinear
		}
		omw := t.OMW(&p)
		fmt.Fprintf(&sb, "%-3d %-15s %5d %3d %3d %3d %3d %7.2f%% %7.2f%%\n",
			i+1, displayName, p.Points, p.Matches, p.Wins, p.Draws, p.Losses,
			p.WinPercentage(), omw)
	}

	// Resumen de partidas
	sb.WriteString("\nResumen de Partidas\n")
	sb.WriteString("--------------------\n")

	for _, match := range t.Matches {
		displayName1 := match.Player1
		displayName2 := match.Player2
		// Añade asterisco al campeón de ese match
		if match.Champion != "" {
			if match.Champion == match.Player1 {
				displayName1 += "*"
			} else if match.Champion == match.Player2 {
				displayName2 += "*"
			}
		}
		fmt.Fprintf(&sb, "%-10s %d - %d  %s\n",
			displayName1, match.GamesWon1, match.GamesWon2, displayName2)
	}

	return sb.String()
}
