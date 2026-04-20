package engine

import (
	"fmt"
	"math"
	"sort"

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

	matchGrid := make([][]*model.Match, rounds+1)
	for r := 1; r <= rounds; r++ {
		size := slots >> uint(r)
		if r == rounds && tournament.ThirdPlaceMatch {
			size = 2
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

	// Sort by seed ascending; seed=0 (unseeded) goes last; tiebreak by CreatedAt.
	sort.SliceStable(registrations, func(i, j int) bool {
		// Treat nil seed pointer as unseeded (0)
		getSeed := func(r *model.TournamentRegistration) int {
			if r.Seed == nil {
				return 0
			}
			return *r.Seed
		}

		si, sj := getSeed(registrations[i]), getSeed(registrations[j])

		if si == 0 && sj != 0 {
			return false
		}
		if si != 0 && sj == 0 {
			return true
		}
		if si != sj {
			return si < sj
		}
		return registrations[i].CreatedAt.Before(registrations[j].CreatedAt)
	})

	// Seed players into round-1 slots using standard tournament seeding
	// so top seeds are spread across the bracket and receive byes.
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

// buildSeededSlots places players into bracket slots following standard
// tournament seeding rules:
//   - Seed 1 and Seed 2 are in opposite halves → can only meet in the Final
//   - Seed 1/2 meet Seed 3/4 only in the Semis, and so on
//   - Byes are awarded to the top seeds (1..byeCount) within their correct
//     bracket section so the seeding structure is preserved
func buildSeededSlots(slots, n int) []int {
	result := make([]int, slots)
	for i := range result {
		result[i] = -1
	}

	byeCount := slots - n
	seedOrder := standardSeedOrder(slots)

	// slots 2k and 2k+1 form match k in round 1.
	// seedOrder[slotIdx] gives the seed number assigned to that slot.
	// For each match pair, the higher seed (lower number) gets the top slot
	// and may receive a bye if their seed number is within byeCount.
	for matchIdx := 0; matchIdx < slots/2; matchIdx++ {
		slotA := matchIdx * 2
		slotB := matchIdx*2 + 1
		seedA := seedOrder[slotA]
		seedB := seedOrder[slotB]

		// Identify which slot holds the higher seed (lower seed number = better rank)
		higherSeed, lowerSeed := seedA, seedB
		higherSlot, lowerSlot := slotA, slotB
		if seedB < seedA {
			higherSeed, lowerSeed = seedB, seedA
			higherSlot, lowerSlot = slotB, slotA
		}

		// Assign the higher-seeded player if they exist
		higherPlayerIdx := higherSeed - 1 // seed is 1-based, playerIdx is 0-based
		if higherPlayerIdx < n {
			result[higherSlot] = higherPlayerIdx
		}

		// Top seeds (seed number within byeCount) receive a bye:
		// leave lowerSlot as -1 so AutoAdvance picks them up.
		// For all other matches, assign the lower-seeded player if they exist.
		if higherSeed <= byeCount {
			// bye — lowerSlot stays -1
		} else {
			lowerPlayerIdx := lowerSeed - 1
			if lowerPlayerIdx < n {
				result[lowerSlot] = lowerPlayerIdx
			}
			// if lowerPlayerIdx >= n it's a natural bye, stays -1
		}
	}

	return result
}

// standardSeedOrder returns a slice of length n where index i contains
// the seed number assigned to slot i, following the standard tournament
// bracket seeding pattern that ensures:
//   - Seed 1 vs Seed 2 only possible in the Final
//   - Seed 1/2 vs Seed 3/4 only possible in the Semis
//   - and so on for each round
//
// Examples:
//
//	n=2: [1, 2]
//	n=4: [1, 4, 3, 2]
//	n=8: [1, 8, 5, 4, 3, 6, 7, 2]
func standardSeedOrder(n int) []int {
	if n == 1 {
		return []int{1}
	}
	prev := standardSeedOrder(n / 2)
	result := make([]int, n)
	for i, s := range prev {
		result[i*2] = s
		result[i*2+1] = n + 1 - s
	}
	return result
}

func nextPowerOf2(n int) int {
	p := 1
	for p < n {
		p <<= 1
	}
	return p
}