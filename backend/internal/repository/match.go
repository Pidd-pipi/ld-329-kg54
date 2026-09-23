package repository

import (
	"cyskillswap/internal/model"
)

// FindMatch returns the match with the given id.
func FindMatch(id int) (model.Match, bool) {
	for _, match := range ListMatches() {
		if match.ID == id {
			return match, true
		}
	}
	return model.Match{}, false
}

// MatchParticipants reports whether user is one of the two parties.
func MatchParticipants(match model.Match) []string {
	return []string{match.Provider, match.Learner}
}
