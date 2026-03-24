// package engine

// import (
// 	"math/rand"
// 	"errors"
// 	"math"

// 	"github.com/uthmanduro/BracketForge/internal/model"
// )

// func nextPowerOfTwo(n int) int {
// 	if n <= 0 {
// 		return 1
// 	}
// 	power := 1
// 	for power < n {
// 		power *= 2
// 	}
// 	return power
// }


// func GenerateSingleEliminationBracket(players []model.Player) ([]model.Match, error) {

// 	if len(players) < 2 {
// 		return nil, errors.New("at least 2 players are required to generate a bracket")
// 	}

// 	// Shuffle players to randomize matchups
// 	rand.Shuffle(len(players), func(i, j int) {
// 		players[i], players[j] = players[j], players[i]
// 	})

// 	// Generate bracket size based on the number of players (next power of 2)
// 	bracketSize := nextPowerOfTwo(len(players))

// 	// Determine the number of rounds needed
// 	totalRounds := int(math.Log2(float64(bracketSize)))

// 	// Create matches for the first round
// 	matches := make([]model.Match, 0, bracketSize-1)
// 	playerIndex := 0

// 	for round := 1; round <= totalRounds; round++ {
// 		matchesInRound := bracketSize / (1 << round)

// 		for i := 0; i < matchesInRound; i++ {
// 			var p1, p2 *model.Player
			
// 			if playerIndex < len(players) {
// 				p1 = &players[playerIndex]
// 				playerIndex++
// 			}

// 			if playerIndex < len(players) {
// 				p2 = &players[playerIndex]
// 				playerIndex++
// 			}

// 			match := model.Match{
// 				Round: round,
// 				Player1: p1,
// 				Player2: p2,
// 			}
// 			matches = append(matches, match)
// 		}

// 	}
	
// 	return matches, nil
// }

package engine

import (
	"fmt"
	"math"

	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"
)

type SingleEliminationEngine struct{ matchRepo *repository.MatchRepo }

func NewSingleEliminationEngine(mr *repository.MatchRepo) *SingleEliminationEngine {
	return &SingleEliminationEngine{matchRepo: mr}
}

func (e *SingleEliminationEngine) GenerateBracket(
	stage *model.Stage,
	tournament *model.Tournament,
	registrations []*model.TournamentRegistration,
) ([]*model.Match, error) {
	n := len(registrations)
	if n < 2 {
		return nil, fmt.Errorf("need at least 2 players, got %d", n)
	}

	bestOf := tournament.BestOf
	if stage.BestOf != nil {
		bestOf = *stage.BestOf
	}

	slots := nextPowerOf2(n)
	rounds := int(math.Log2(float64(slots)))

	// matchGrid[round][position]: round 1 = first round played, round=rounds = final.
	matchGrid := make([][]*model.Match, rounds+1)
	for r := 1; r <= rounds; r++ {
		size := slots >> uint(r)
		if r == rounds && tournament.ThirdPlaceMatch {
			size = 2 // index 0 = final, index 1 = third-place
		}
		matchGrid[r] = make([]*model.Match, size)
	}

	// Insert match rows final-first so next_match_id references exist.
	for r := rounds; r >= 1; r-- {
		for pos := 0; pos < len(matchGrid[r]); pos++ {
			isThirdPlace := r == rounds && pos == 1 && tournament.ThirdPlaceMatch

			round := r
			position := pos + 1
			m := &model.Match{
				StageID:       stage.ID,
				Round:         &round,
				MatchPosition: &position,
				BestOf:        bestOf,
				Status:        "pending",
			}

			if !isThirdPlace && r < rounds {
				parentPos := pos / 2
				parent := matchGrid[r+1][parentPos]
				parentID := parent.ID
				nextSlot := (pos % 2) + 1
				m.NextMatchID = &parentID
				m.NextMatchSlot = &nextSlot

				if tournament.ThirdPlaceMatch && r == rounds-1 {
					tpID := matchGrid[rounds][1].ID
					loserSlot := pos + 1
					m.LoserNextMatchID = &tpID
					m.LoserNextMatchSlot = &loserSlot
				}
			}

			if err := e.matchRepo.Create(m); err != nil {
				return nil, fmt.Errorf("create match r%d p%d: %w", r, pos, err)
			}
			matchGrid[r][pos] = m
		}
	}

	// Seed players into round-1 slots.
	seededSlots := buildSeededSlots(slots, n)
	round1 := matchGrid[1]

	for slotIdx, regIdx := range seededSlots {
		if regIdx < 0 {
			continue
		}
		matchIdx := slotIdx / 2
		slot := (slotIdx % 2) + 1
		p := &model.MatchParticipant{
			MatchID:        round1[matchIdx].ID,
			RegistrationID: registrations[regIdx].ID,
			Slot:           slot,
		}
		if err := e.matchRepo.AddParticipant(p); err != nil {
			return nil, fmt.Errorf("add participant slot %d: %w", slotIdx, err)
		}
	}

	// Auto-complete BYE matches.
	byeEng := NewByeEngine(e.matchRepo)
	for _, m := range round1 {
		parts, err := e.matchRepo.GetParticipants(m.ID)
		if err != nil {
			return nil, err
		}
		if len(parts) == 1 {
			if err := byeEng.AutoAdvance(m, parts[0]); err != nil {
				return nil, fmt.Errorf("bye advance %s: %w", m.ID, err)
			}
		}
	}

	var all []*model.Match
	for r := 1; r <= rounds; r++ {
		all = append(all, matchGrid[r]...)
	}
	return all, nil
}

func buildSeededSlots(slots, n int) []int {
	result := make([]int, slots)
	for i := range result {
		result[i] = -1
	}
	positions := bracketPositions(slots)
	for i, pos := range positions {
		if i < n {
			result[pos] = i
		}
	}
	return result
}

func bracketPositions(size int) []int {
	pos := []int{0}
	for len(pos) < size {
		next := make([]int, 0, len(pos)*2)
		for _, p := range pos {
			next = append(next, p, size-1-p)
		}
		pos = next
	}
	return pos
}

func nextPowerOf2(n int) int {
	p := 1
	for p < n {
		p <<= 1
	}
	return p
}