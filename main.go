package pi

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/arewabolu/csvmanager"
	"golang.org/x/exp/slices"
)

// Calculate New Home rating for both teams
// TODO: make a function for creating a new rating file
// func needs the team names that's all

const (
	delta  = 2.5
	gamma  = 0.79
	lambda = 0.054
	Mu     = 0.01
	phi    = 1 //form factor
)

type Team struct {
	Name                      string
	HomeRating                float64
	AwayRating                float64
	ContinuousPerformanceHome int
	ContinuousPerformanceAway int //continuous performance is either over or under
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

func (t *Team) ProvisionalRatingHome() float64 {
	sub := t.ContinuousPerformanceHome - 1
	f := float64(sub)
	denum := math.Pow(f, delta)
	total := f / denum
	return t.HomeRating + (Mu * total)
}

func (t *Team) ProvisionalRatingAway() float64 {
	sub := t.ContinuousPerformanceAway - 1
	f := float64(sub)
	denum := -math.Pow(math.Abs(f), delta)
	total := f / denum
	return t.AwayRating + (-Mu * total)
}

func (t *Team) ProvisionalRatingAwayV2() float64 {
	sub := t.ContinuousPerformanceAway - 1
	f := float64(sub)
	denum := math.Pow(f, delta)
	total := f / denum
	return t.AwayRating + (Mu * total)
}

func (t *Team) ProvisionalRatingHomeV2() float64 {
	sub := t.ContinuousPerformanceHome - 1
	f := float64(sub)
	denum := -math.Pow(math.Abs(f), delta)
	total := f / denum
	return t.HomeRating + (-Mu * total)
}

// returns home and away background ratings for a given team
func (t *Team) UpdateBackgroundHomeTeamRatings(errorGDFunc float64) {
	BRH := t.HomeRating + (errorGDFunc * lambda)
	BRA := t.AwayRating + ((BRH - t.HomeRating) * gamma)
	t.HomeRating = BRH
	t.AwayRating = BRA
}

func (t *Team) UpdateBackgroundAwayTeamRatings(errorGDFunc float64) {
	BRA := t.AwayRating + (errorGDFunc * lambda)
	BRH := t.HomeRating + ((BRA - t.AwayRating) * gamma)
	t.AwayRating = BRA
	t.HomeRating = BRH
}

func (t *Team) ResetContinuousPerformanceHome() {
	t.ContinuousPerformanceHome = 0
}

func (t *Team) ResetContinuousPerformanceAway() {
	t.ContinuousPerformanceAway = 0
}

func (t *Team) UpdateContinuousPerformanceHome() {
	t.ContinuousPerformanceHome = t.ContinuousPerformanceHome + 1
}

func (t *Team) UpdateContinuousPerformanceAway() {
	t.ContinuousPerformanceAway = t.ContinuousPerformanceAway - 1
}

func (t *Team) UpdateContinuousPerformanceHomeV2() {
	t.ContinuousPerformanceHome = t.ContinuousPerformanceHome - 1
}

func (t *Team) UpdateContinuousPerformanceAwayV2() {
	t.ContinuousPerformanceAway = t.ContinuousPerformanceAway + 1
}

func UpdateTeamRatings(filepath string, homeTeamName, awayTeamName string, homeGoalScored, awayGoalScored int) error {
	ratings, err := csvmanager.ReadCsv(filepath, 0755, true)
	if err != nil {
		panic(err)
	}
	TeamCol := ratings.Col("TeamName").String()

	HomeTeam := &Team{Name: homeTeamName}
	AwayTeam := &Team{Name: awayTeamName}
	HToccurence := slices.Index(TeamCol, HomeTeam.Name)
	if HToccurence == -1 {
		return errors.New("couldn't find specified team " + HomeTeam.Name)
	}
	AToccurence := slices.Index(TeamCol, AwayTeam.Name)
	if AToccurence == -1 {
		return errors.New("couldn't find specified team " + AwayTeam.Name)
	}
	if HToccurence > -1 {
		HTData := ratings.Row(HToccurence).String()
		HomeTeam.HomeRating, err = strconv.ParseFloat(HTData[1], 64)
		if err != nil {
			panic(err)
		}
		HomeTeam.AwayRating, err = strconv.ParseFloat(HTData[2], 64)
		if err != nil {
			panic(err)
		}
		HomeTeam.ContinuousPerformanceHome, err = strconv.Atoi(HTData[3])
		if err != nil {
			panic(err)
		}
		HomeTeam.ContinuousPerformanceAway, err = strconv.Atoi(HTData[4])
		if err != nil {
			panic(err)
		}
	}
	if AToccurence > -1 {
		ATData := ratings.Row(AToccurence).String()
		AwayTeam.HomeRating, err = strconv.ParseFloat(ATData[1], 64)
		if err != nil {
			panic(err)
		}
		AwayTeam.AwayRating, err = strconv.ParseFloat(ATData[2], 64)
		if err != nil {
			panic(err)
		}
		AwayTeam.ContinuousPerformanceHome, err = strconv.Atoi(ATData[3])
		if err != nil {
			panic(err)
		}
		AwayTeam.ContinuousPerformanceAway, err = strconv.Atoi(ATData[4])
		if err != nil {
			panic(err)
		}
	}
	HxG := ExpectedGoalIndividual(HomeTeam.HomeRating)
	AxG := ExpectedGoalIndividual(AwayTeam.AwayRating)
	xGD := ExpectedGoalDifference(HxG, AxG)
	GD := goalDifference(homeGoalScored, awayGoalScored)
	var HomeErrFunc, AwayErrFunc float64
	errFunc := errorGDFunc(errorGD(GD, xGD))
	if xGD > float64(GD) {
		HomeErrFunc = -errFunc
		AwayErrFunc = errFunc

	} else {
		HomeErrFunc = errFunc
		AwayErrFunc = -errFunc
	}
	HomeTeam.UpdateBackgroundHomeTeamRatings(HomeErrFunc)
	AwayTeam.UpdateBackgroundAwayTeamRatings(AwayErrFunc)

	switch {
	case xGD >= 0 && GD > 0:
		HomeTeam.UpdateContinuousPerformanceHome()
		AwayTeam.UpdateContinuousPerformanceAway()
	case xGD > 0 && GD < 0:
		HomeTeam.ResetContinuousPerformanceHome()
		HomeTeam.UpdateContinuousPerformanceHomeV2()
		AwayTeam.ResetContinuousPerformanceAway()
		AwayTeam.UpdateContinuousPerformanceAwayV2()
	case xGD < 0 && GD > 0:
		AwayTeam.ResetContinuousPerformanceAway()
		AwayTeam.UpdateContinuousPerformanceAway()
		HomeTeam.ResetContinuousPerformanceHome()
		HomeTeam.UpdateContinuousPerformanceHome()
	case xGD <= 0 && GD < 0:
		AwayTeam.UpdateContinuousPerformanceAwayV2()
		HomeTeam.UpdateContinuousPerformanceHomeV2()
	case xGD > 0 || xGD < 0 && GD == 0:
		HomeTeam.ResetContinuousPerformanceHome()
		AwayTeam.ResetContinuousPerformanceAway()
	}
	HTWrite := []string{HomeTeam.Name, fmt.Sprintf("%.2f", HomeTeam.HomeRating), fmt.Sprintf("%.2f", HomeTeam.AwayRating), fmt.Sprintf("%d", HomeTeam.ContinuousPerformanceHome), fmt.Sprintf("%d", HomeTeam.ContinuousPerformanceAway)}
	ATWrite := []string{AwayTeam.Name, fmt.Sprintf("%.2f", AwayTeam.HomeRating), fmt.Sprintf("%.2f", AwayTeam.AwayRating), fmt.Sprintf("%d", AwayTeam.ContinuousPerformanceHome), fmt.Sprintf("%d", AwayTeam.ContinuousPerformanceAway)}
	csvmanager.ReplaceRow(filepath, 0755, HToccurence+1, HTWrite)
	csvmanager.ReplaceRow(filepath, 0755, AToccurence+1, ATWrite)
	return nil
}

func CheckRatings(filepath string, team string) []string {
	ratings, err := csvmanager.ReadCsv(filepath, 0755, true)
	if err != nil {
		panic(err)
	}
	TeamCol := ratings.Col("TeamName").String()
	index := slices.Index(TeamCol, team)
	if index != -1 {
		return ratings.Row(index).String()
	}
	return []string{}
}

func Search(filepath string, teamName, venue string) *Team {
	teamInfo := CheckRatings(filepath, teamName)
	if len(teamInfo) == 0 {
		return &Team{}
	}
	team := &Team{}
	team.Name = teamName
	var err error
	team.HomeRating, err = strconv.ParseFloat(teamInfo[1], 64)
	if err != nil {
		panic(err)
	}
	team.AwayRating, err = strconv.ParseFloat(teamInfo[2], 64)
	if err != nil {
		panic(err)
	}
	team.ContinuousPerformanceHome, err = strconv.Atoi(teamInfo[3])
	if err != nil {
		panic(err)
	}
	team.ContinuousPerformanceAway, err = strconv.Atoi(teamInfo[4])
	if err != nil {
		panic(err)
	}
	if venue == "away" {
		if team.ContinuousPerformanceAway > 1 {
			team.AwayRating = team.ProvisionalRatingAwayV2()
		} else if team.ContinuousPerformanceAway < -1 {
			team.AwayRating = team.ProvisionalRatingAway()
		}
	}
	if venue == "home" {
		if team.ContinuousPerformanceHome > 1 {
			team.HomeRating = team.ProvisionalRatingHome()
		} else if team.ContinuousPerformanceHome < -1 {
			team.HomeRating = team.ProvisionalRatingHomeV2()
		}
	}
	return team
}

func TotalBackgroundRating(homeRating, awayRating float64) float64 {
	return homeRating + awayRating
}
