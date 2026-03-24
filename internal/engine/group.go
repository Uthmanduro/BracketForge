// package engine

// import (
// 	"errors"
// 	"math/rand"

// 	"github.com/uthmanduro/BracketForge/internal/model"
// )

// func DistributePlayersIntoGroups(players []model.Player, numGroups int) ([]model.Group, error) {
// 	if numGroups <= 0 {
// 		return nil, errors.New("number of groups must be positive")
// 	}

// 	// Shuffle players to randomize group distribution
// 	rand.Shuffle(len(players), func(i, j int) {
// 		players[i], players[j] = players[j], players[i]
// 	})

// 	groups := make([]model.Group, numGroups)
// 	for i := range groups {
// 		groups[i].Name = "Group " + string('A'+i)

// 	}

// 	for i, player := range players {
// 		groupIndex := i % numGroups
// 		groups[groupIndex].Players = append(groups[groupIndex].Players, player)
// 	}

// 	return groups, nil
// }


package engine

import (
	"fmt"
	"sort"

	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"
)

type GroupEngine struct {
	stageRepo     *repository.StageRepo
	matchRepo     *repository.MatchRepo
	standingsRepo *repository.StandingsRepo
	rrEngine      *RoundRobinEngine
}

func NewGroupEngine(
	sr *repository.StageRepo,
	mr *repository.MatchRepo,
	standingsRepo *repository.StandingsRepo,
	rrEngine *RoundRobinEngine,
) *GroupEngine {
	return &GroupEngine{stageRepo: sr, matchRepo: mr, standingsRepo: standingsRepo, rrEngine: rrEngine}
}

func (e *GroupEngine) Draw(
	stage *model.Stage,
	tournament *model.Tournament,
	registrations []*model.TournamentRegistration,
	numberOfGroups int,
) ([]*model.Group, error) {
	if numberOfGroups < 2 {
		return nil, fmt.Errorf("need at least 2 groups")
	}
	if len(registrations) < numberOfGroups {
		return nil, fmt.Errorf("not enough players for %d groups", numberOfGroups)
	}

	groups := make([]*model.Group, numberOfGroups)
	for i := 0; i < numberOfGroups; i++ {
		g := &model.Group{
			StageID: stage.ID,
			Name:    fmt.Sprintf("Group %s", string(rune('A'+i))),
		}
		if err := e.stageRepo.CreateGroup(g); err != nil {
			return nil, fmt.Errorf("create group %s: %w", g.Name, err)
		}
		groups[i] = g
	}

	// Serpentine seeding.
	groupAssignments := make([][]*model.TournamentRegistration, numberOfGroups)
	for i, reg := range registrations {
		row := i / numberOfGroups
		col := i % numberOfGroups
		if row%2 == 1 {
			col = numberOfGroups - 1 - col
		}
		groupAssignments[col] = append(groupAssignments[col], reg)
	}

	for gi, group := range groups {
		for _, reg := range groupAssignments[gi] {
			gr := &model.GroupRegistration{GroupID: group.ID, RegistrationID: reg.ID}
			if err := e.stageRepo.AddGroupRegistration(gr); err != nil {
				return nil, fmt.Errorf("assign player to group: %w", err)
			}
		}
		gid := group.ID
		if _, err := e.rrEngine.GenerateMatches(stage, tournament, groupAssignments[gi], &gid); err != nil {
			return nil, fmt.Errorf("generate matches for %s: %w", group.Name, err)
		}
	}
	return groups, nil
}

// RankedStandings returns standings for a group with tiebreakers applied.
func (e *GroupEngine) RankedStandings(groupID string) ([]*model.Standings, error) {
	standings, err := e.standingsRepo.ListByGroup(groupID)
	if err != nil {
		return nil, err
	}
	if len(standings) <= 1 {
		for i, s := range standings {
			s.Rank = i + 1
		}
		return standings, nil
	}
	sorted := applyTiebreaker(standings, groupID, e.matchRepo)
	for i, s := range sorted {
		s.Rank = i + 1
	}
	return sorted, nil
}

// ── Tiebreaker cascade ─────────────────────────────────────────────────────

func applyTiebreaker(standings []*model.Standings, groupID string, matchRepo *repository.MatchRepo) []*model.Standings {
	groups := groupByPoints(standings)
	var result []*model.Standings
	for _, group := range groups {
		if len(group) == 1 {
			result = append(result, group...)
			continue
		}
		result = append(result, resolveTied(group, groupID, matchRepo)...)
	}
	return result
}

func groupByPoints(standings []*model.Standings) [][]*model.Standings {
	sort.Slice(standings, func(i, j int) bool {
		return standings[i].Points > standings[j].Points
	})
	var groups [][]*model.Standings
	i := 0
	for i < len(standings) {
		j := i + 1
		for j < len(standings) && standings[j].Points == standings[i].Points {
			j++
		}
		groups = append(groups, standings[i:j])
		i = j
	}
	return groups
}

func resolveTied(tied []*model.Standings, groupID string, matchRepo *repository.MatchRepo) []*model.Standings {
	if len(tied) <= 1 {
		return tied
	}
	ids := make([]string, len(tied))
	for i, s := range tied {
		ids[i] = s.RegistrationID
	}

	// Step 2: head-to-head points among tied players.
	h2h := computeH2HStandings(tied, groupID, ids, matchRepo)
	if h2h != nil && !allPointsEqual(h2h) {
		return applyTiebreaker(h2h, groupID, matchRepo)
	}

	// Steps 3–4: mutual-match set/game ratios.
	mutualStats, err := matchRepo.GetMutualSetStats(groupID, ids)
	if err == nil {
		// Step 3: mutual set ratio.
		byMutualSet := cloneStandings(tied)
		sort.SliceStable(byMutualSet, func(i, j int) bool {
			a := mutualStats[byMutualSet[i].RegistrationID]
			b := mutualStats[byMutualSet[j].RegistrationID]
			return ratio(a.SetsWon, a.SetsLost) > ratio(b.SetsWon, b.SetsLost)
		})
		if !allMutualEqual(byMutualSet, mutualStats, func(s *model.MutualSetStats) float64 {
			return ratio(s.SetsWon, s.SetsLost)
		}) {
			return byMutualSet
		}

		// Step 4: mutual game ratio.
		byMutualGame := cloneStandings(tied)
		sort.SliceStable(byMutualGame, func(i, j int) bool {
			a := mutualStats[byMutualGame[i].RegistrationID]
			b := mutualStats[byMutualGame[j].RegistrationID]
			return ratio(a.GamesWon, a.GamesLost) > ratio(b.GamesWon, b.GamesLost)
		})
		if !allMutualEqual(byMutualGame, mutualStats, func(s *model.MutualSetStats) float64 {
			return ratio(s.GamesWon, s.GamesLost)
		}) {
			return byMutualGame
		}
	}

	// Step 5: overall set ratio.
	bySet := cloneStandings(tied)
	sort.SliceStable(bySet, func(i, j int) bool {
		return ratio(bySet[i].SetsWon, bySet[i].SetsLost) > ratio(bySet[j].SetsWon, bySet[j].SetsLost)
	})
	if !allEqual(bySet, func(s *model.Standings) float64 { return ratio(s.SetsWon, s.SetsLost) }) {
		return bySet
	}

	// Step 6: overall game ratio.
	byGame := cloneStandings(tied)
	sort.SliceStable(byGame, func(i, j int) bool {
		return ratio(byGame[i].GamesWon, byGame[i].GamesLost) > ratio(byGame[j].GamesWon, byGame[j].GamesLost)
	})
	if !allEqual(byGame, func(s *model.Standings) float64 { return ratio(s.GamesWon, s.GamesLost) }) {
		return byGame
	}

	// Step 7: draw — return as-is.
	return tied
}

func computeH2HStandings(
	tied []*model.Standings,
	groupID string,
	ids []string,
	matchRepo *repository.MatchRepo,
) []*model.Standings {
	mutualMatches, err := matchRepo.GetMutualMatches(groupID, ids)
	if err != nil || len(mutualMatches) == 0 {
		return nil
	}

	h2h := make(map[string]*model.Standings, len(tied))
	for _, s := range tied {
		h2h[s.RegistrationID] = &model.Standings{
			RegistrationID: s.RegistrationID,
			PlayerName:     s.PlayerName,
			StageID:        s.StageID,
			GroupID:        s.GroupID,
		}
	}

	for _, m := range mutualMatches {
		if m.WinnerRegistrationID == nil {
			continue
		}
		winnerID := *m.WinnerRegistrationID
		parts, err := matchRepo.GetParticipants(m.ID)
		if err != nil {
			continue
		}
		for _, p := range parts {
			if s, ok := h2h[p.RegistrationID]; ok {
				s.Played++
				if p.RegistrationID != winnerID {
					s.Losses++
				}
			}
		}
		if s, ok := h2h[winnerID]; ok {
			s.Wins++
			s.Points += 2
		}
	}

	result := make([]*model.Standings, 0, len(h2h))
	for _, s := range h2h {
		result = append(result, s)
	}
	return result
}

func allPointsEqual(standings []*model.Standings) bool {
	if len(standings) == 0 {
		return true
	}
	p := standings[0].Points
	for _, s := range standings[1:] {
		if s.Points != p {
			return false
		}
	}
	return true
}

func ratio(won, lost int) float64 {
	if lost == 0 {
		if won == 0 {
			return 0
		}
		return float64(won) * 1000
	}
	return float64(won) / float64(lost)
}

func allEqual(s []*model.Standings, fn func(*model.Standings) float64) bool {
	if len(s) == 0 {
		return true
	}
	v := fn(s[0])
	for _, item := range s[1:] {
		if fn(item) != v {
			return false
		}
	}
	return true
}

func allMutualEqual(s []*model.Standings, stats map[string]*model.MutualSetStats, fn func(*model.MutualSetStats) float64) bool {
	if len(s) == 0 {
		return true
	}
	v := fn(stats[s[0].RegistrationID])
	for _, item := range s[1:] {
		if fn(stats[item.RegistrationID]) != v {
			return false
		}
	}
	return true
}

func cloneStandings(src []*model.Standings) []*model.Standings {
	out := make([]*model.Standings, len(src))
	copy(out, src)
	return out
}