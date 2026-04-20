package engine

import (
	"fmt"

	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"
	"gorm.io/gorm"
)

type RoundRobinEngine struct {
	matchRepo     *repository.MatchRepo
	standingsRepo *repository.StandingsRepo
}

func NewRoundRobinEngine(mr *repository.MatchRepo, sr *repository.StandingsRepo) *RoundRobinEngine {
	return &RoundRobinEngine{matchRepo: mr, standingsRepo: sr}
}

func (e *RoundRobinEngine) GenerateMatches(
	stage *model.Stage,
	tournament *model.Tournament,
	registrations []*model.TournamentRegistration,
	groupID *string,
) ([]*model.Match, error) {
	n := len(registrations)
	if n < 2 {
		// return nil, fmt.Errorf("need at least 2 players, got %d", n)
	}

	bestOf := tournament.BestOf
	if stage.BestOf != nil {
		bestOf = *stage.BestOf
	}

	// Pad to even number with a nil ghost for the rotation algorithm.
	effective := n
	if n%2 != 0 {
		effective = n + 1
	}

	numRounds := effective - 1
	matchesPerRound := effective / 2

	players := make([]*model.TournamentRegistration, n)
	copy(players, registrations)

	// Circle algorithm: fix players[0], rotate the rest.
	circle := make([]*model.TournamentRegistration, effective-1)
	for i := 0; i < effective-1; i++ {
		if i < n-1 {
			circle[i] = players[i+1]
		}
		// nil = ghost BYE slot (odd player count)
	}
	fixed := players[0]

	var allMatches []*model.Match

	for round := 1; round <= numRounds; round++ {
		pos := 0
		for i := 0; i < matchesPerRound; i++ {
			var p1, p2 *model.TournamentRegistration
			if i == 0 {
				p1 = fixed
				p2 = circle[effective-2]
			} else {
				p1 = circle[i-1]
				p2 = circle[effective-2-i]
			}

			// Skip ghost BYE pairings.
			if p1 == nil || p2 == nil {
				continue
			}

			pos++
			r, p := round, pos
			m := &model.Match{
				StageID:       stage.ID,
				GroupID:       groupID,
				Round:         &r,
				MatchPosition: &p,
				BestOf:        bestOf,
				Status:        "pending",
			}

			if err := e.matchRepo.Create(m); err != nil {
				return nil, fmt.Errorf("create rr match r%d: %w", round, err)
			}

			part1 := &model.MatchParticipant{MatchID: m.ID, RegistrationID: p1.ID, Slot: 1}
			part2 := &model.MatchParticipant{MatchID: m.ID, RegistrationID: p2.ID, Slot: 2}
			if err := e.matchRepo.AddParticipant(part1); err != nil {
				return nil, err
			}
			if err := e.matchRepo.AddParticipant(part2); err != nil {
				return nil, err
			}
			allMatches = append(allMatches, m)
		}

		// Rotate circle left by 1.
		last := circle[len(circle)-1]
		copy(circle[1:], circle[:len(circle)-1])
		circle[0] = last
	}

	// Initialise standings rows for every player.
	for _, reg := range registrations {
		if err := e.standingsRepo.InitForRegistration(stage.ID, groupID, reg.ID); err != nil {
			return nil, fmt.Errorf("init standings: %w", err)
		}
	}

	return allMatches, nil
}

func (e *RoundRobinEngine) WithTx(tx *gorm.DB) *RoundRobinEngine {
	return &RoundRobinEngine{
		matchRepo:     e.matchRepo.WithTx(tx),
		standingsRepo: e.standingsRepo.WithTx(tx),
	}
}