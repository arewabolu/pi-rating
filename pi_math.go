package pi

import (
	"math"
	"strings"
)

const (
	delta  = 2.5
	gamma  = 0.79  //0.70
	lambda = 0.054 //  0.035/
	Mu     = 0.01
	phi    = 1 //form factor
)

type Team struct {
	Name       string
	HomeRating float64
	AwayRating float64
	// Continuous Performance values are a measure of a team's recent form at home
	ContinuousPerformanceHome int
	// Continuous Performance values are a measure of a team's recent form at away
	ContinuousPerformanceAway int
	index                     int
}

func ExpectedGoalIndividual(rating float64) float64 {
	RG := math.Abs(rating) / 3
	if rating < 0 {
		return -1 * (math.Pow(10, RG) - 1)
	}
	return math.Pow(10, RG) - 1
}

func ExpectedGoalDifference(homexG, awayxG float64) float64 { return homexG - awayxG }

func errorGD(goalDifference int, expectedGoalDifference float64) float64 {
	return math.Abs(float64(goalDifference) - expectedGoalDifference)
}

func errorGDFunc(errorGD float64) float64 { return 3 * math.Log10(1+errorGD) }

func goalDifference(homeGoal, awayGoal int) int { return homeGoal - awayGoal }

func (t Team) provisionalRatingHome() float64 {
	sub := t.ContinuousPerformanceHome - 1
	f := float64(sub)
	denum := math.Pow(f, delta)
	total := f / denum
	return t.HomeRating + (Mu * total)
}

func (t Team) provisionalRatingAway() float64 {
	sub := t.ContinuousPerformanceAway - 1
	f := float64(sub)
	denum := math.Pow(math.Abs(f), delta)
	total := f / denum
	return t.AwayRating + (-Mu * total)
}

func (t Team) provisionalRatingAwayV2() float64 {
	sub := t.ContinuousPerformanceAway - 1
	f := float64(sub)
	denum := math.Pow(f, delta)
	total := f / denum
	return t.AwayRating + (Mu * total)
}

func (t Team) provisionalRatingHomeV2() float64 {
	sub := t.ContinuousPerformanceHome - 1
	f := float64(sub)
	denum := math.Pow(math.Abs(f), delta)
	total := f / denum
	return t.HomeRating + (-Mu * total)
}

// Should be used to incorporate form into the team ratings
func (t Team) ProvisionalRating(venue string) Team {
	venue = strings.ToLower(venue)
	switch venue {
	case "away":
		if t.ContinuousPerformanceAway > 1 {
			t.AwayRating = t.provisionalRatingAwayV2()
		} else if t.ContinuousPerformanceAway < -1 {
			t.AwayRating = t.provisionalRatingAway()
		}
	case "home":
		if t.ContinuousPerformanceHome > 1 {
			t.HomeRating = t.provisionalRatingHome()
		} else if t.ContinuousPerformanceHome < -1 {
			t.HomeRating = t.provisionalRatingHomeV2()
		}
	default:
		return t
	}

	return t
}

// revises home and away background ratings for a given team
func (t *Team) updateBackgroundHometeamRatings(errorGDFunc float64) {
	BRH := t.HomeRating + (errorGDFunc * lambda)
	BRA := t.AwayRating + ((BRH - t.HomeRating) * gamma)
	t.HomeRating = BRH
	t.AwayRating = BRA
}

// revises home and away background ratings for a given team
func (t *Team) updateBackgroundAwayteamRatings(errorGDFunc float64) {
	BRA := t.AwayRating + (errorGDFunc * lambda)
	BRH := t.HomeRating + ((BRA - t.AwayRating) * gamma)
	t.AwayRating = BRA
	t.HomeRating = BRH
}

func (t *Team) resetContinuousPerformanceHome() {
	t.ContinuousPerformanceHome = 0
}

func (t *Team) resetContinuousPerformanceAway() {
	t.ContinuousPerformanceAway = 0
}

func (t *Team) updateContinuousPerformanceHome() {
	t.ContinuousPerformanceHome = t.ContinuousPerformanceHome + 1
}

func (t *Team) updateContinuousPerformanceAway() {
	t.ContinuousPerformanceAway = t.ContinuousPerformanceAway - 1
}

func (t *Team) updateContinuousPerformanceHomeV2() {
	t.ContinuousPerformanceHome = t.ContinuousPerformanceHome - 1
}

func (t *Team) updateContinuousPerformanceAwayV2() {
	t.ContinuousPerformanceAway = t.ContinuousPerformanceAway + 1
}

func RatingDifference(homeRating, awayRating float64) float64 {
	return homeRating - awayRating
}
