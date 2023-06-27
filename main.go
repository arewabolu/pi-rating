package pi

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/arewabolu/csvmanager"
	"golang.org/x/exp/slices"
)

// Calculate New Home rating for both teams
// TODO: make a function for creating a new rating file
// func needs the team names that's all
// TODO: handle all errors properly removing panic,
//
//	provide another way to manage team form
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

func (t *Team) provisionalRatingHome() float64 {
	sub := t.ContinuousPerformanceHome - 1
	f := float64(sub)
	denum := math.Pow(f, delta)
	total := f / denum
	return t.HomeRating + (Mu * total)
}

func (t *Team) provisionalRatingAway() float64 {
	sub := t.ContinuousPerformanceAway - 1
	f := float64(sub)
	denum := math.Pow(math.Abs(f), delta)
	total := f / denum
	return t.AwayRating + (-Mu * total)
}

func (t *Team) provisionalRatingAwayV2() float64 {
	sub := t.ContinuousPerformanceAway - 1
	f := float64(sub)
	denum := math.Pow(f, delta)
	total := f / denum
	return t.AwayRating + (Mu * total)
}

func (t *Team) provisionalRatingHomeV2() float64 {
	sub := t.ContinuousPerformanceHome - 1
	f := float64(sub)
	denum := math.Pow(math.Abs(f), delta)
	total := f / denum
	return t.HomeRating + (-Mu * total)
}

// revises home and away background ratings for a given team
func (t *Team) updateBackgroundHomeTeamRatings(errorGDFunc float64) {
	BRH := t.HomeRating + (errorGDFunc * lambda)
	BRA := t.AwayRating + ((BRH - t.HomeRating) * gamma)
	t.HomeRating = BRH
	t.AwayRating = BRA
}

// revises home and away background ratings for a given team
func (t *Team) updateBackgroundAwayTeamRatings(errorGDFunc float64) {
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

func SetTeamInfo(filepath string, teamName string) (Team, error) {
	teamData, err := CheckRatings(filepath, teamName)
	if err != nil {
		return Team{}, err
	}
	team := Team{Name: teamName}
	team.HomeRating, err = strconv.ParseFloat(teamData[1], 64)
	if err != nil {
		panic(err)
	}
	team.AwayRating, err = strconv.ParseFloat(teamData[2], 64)
	if err != nil {
		panic(err)
	}
	team.ContinuousPerformanceHome, err = strconv.Atoi(teamData[3])
	if err != nil {
		panic(err)
	}
	team.ContinuousPerformanceAway, err = strconv.Atoi(teamData[4])
	if err != nil {
		panic(err)
	}
	team.index, err = strconv.Atoi(teamData[5])
	if err != nil {
		panic(err)
	}
	return team, nil
}

// Update total team information in the Team struct
//
// Error should be checked first before trying to use team struct
func UpdateTeamRatings(filepath string, homeTeamName, awayTeamName string, homeGoalScored, awayGoalScored int) (Team, Team, error) {
	HTData, err := SetTeamInfo(filepath, homeTeamName)
	if err != nil {
		return Team{}, Team{}, err
	}
	ATData, err := SetTeamInfo(filepath, awayTeamName)
	if err != nil {
		return Team{}, Team{}, err
	}

	HxG := ExpectedGoalIndividual(HTData.HomeRating)
	AxG := ExpectedGoalIndividual(ATData.AwayRating)
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
	HTData.updateBackgroundHomeTeamRatings(HomeErrFunc)
	ATData.updateBackgroundAwayTeamRatings(AwayErrFunc)

	switch {
	case xGD >= 0 && GD > 0:
		HTData.updateContinuousPerformanceHome()
		ATData.updateContinuousPerformanceAway()
	case xGD > 0 && GD < 0:
		HTData.resetContinuousPerformanceHome()
		HTData.updateContinuousPerformanceHomeV2()
		ATData.resetContinuousPerformanceAway()
		ATData.updateContinuousPerformanceAwayV2()
	case xGD < 0 && GD > 0:
		ATData.resetContinuousPerformanceAway()
		ATData.updateContinuousPerformanceAway()
		HTData.resetContinuousPerformanceHome()
		HTData.updateContinuousPerformanceHome()
	case xGD <= 0 && GD < 0:
		ATData.updateContinuousPerformanceAwayV2()
		HTData.updateContinuousPerformanceHomeV2()
	case xGD > 0 || xGD < 0 && GD == 0:
		HTData.resetContinuousPerformanceHome()
		ATData.resetContinuousPerformanceAway()
	}

	return HTData, ATData, nil
}

// Write updated rating to rating file
func (t *Team) WriteRatings(filepath string) error {
	team := []string{t.Name, fmt.Sprintf("%.2f", t.HomeRating), fmt.Sprintf("%.2f", t.AwayRating), fmt.Sprintf("%d", t.ContinuousPerformanceHome), fmt.Sprintf("%d", t.ContinuousPerformanceAway)}
	err := csvmanager.ReplaceRow(filepath, 0755, t.index+1, team)
	return err
}

// Checks if a teams is listed in the rating table
// and returns the team information if it exists.
func CheckRatings(filepath string, team string) ([]string, error) {
	ratings, err := csvmanager.ReadCsv(filepath, 0755, true)
	if err != nil {
		return nil, fmt.Errorf("%s not found", filepath)
	}

	if !slices.Contains(ratings.ListHeaders(), "TeamName") {
		return nil, err
	}
	TeamCol := ratings.Col("TeamName").String()
	index := slices.Index(TeamCol, team)
	if index != -1 {
		occurence := ratings.Row(index).String()
		occurence = append(occurence, fmt.Sprint(index))
		return occurence, nil
	}
	return nil, fmt.Errorf("%s not registered. please make sure team is already registered", team)
}

// Search for a teams rating
//
// Error should be checked first before trying to use team struct
func Search(filepath string, teamName string) (Team, error) {
	teamInfo, err := CheckRatings(filepath, teamName)
	if err != nil {
		return Team{}, err
	}
	team := Team{}
	team.Name = teamName
	team.HomeRating, err = strconv.ParseFloat(teamInfo[1], 64)
	if err != nil {
		return Team{}, fmt.Errorf("%s: %v is not a number", teamName, team.HomeRating)
	}
	team.AwayRating, err = strconv.ParseFloat(teamInfo[2], 64)
	if err != nil {
		return Team{}, fmt.Errorf("%s: %v is not a number", teamName, team.AwayRating)
	}
	team.ContinuousPerformanceHome, err = strconv.Atoi(teamInfo[3])
	if err != nil {
		return Team{}, fmt.Errorf("%s: %v is not a number", teamName, team.ContinuousPerformanceHome)
	}
	team.ContinuousPerformanceAway, err = strconv.Atoi(teamInfo[4])
	if err != nil {
		return Team{}, fmt.Errorf("%s: %v is not a number", teamName, team.ContinuousPerformanceAway)

	}

	return team, nil
}

func RatingDifference(homeRating, awayRating float64) float64 {
	return homeRating - awayRating
}

// Should be used to incorporate form into the team ratins
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
