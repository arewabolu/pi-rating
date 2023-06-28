package pi

import (
	"fmt"
	"strconv"
)

// Update total team information in the team struct
//
// Error should be checked first before trying to use team struct
func UpdateTeamRatings(filepath string, hometeamName, awayteamName string, homeGoalScored, awayGoalScored int) (Team, Team, error) {
	HTData, err := setteamInfo(filepath, hometeamName)
	if err != nil {
		return Team{}, Team{}, err
	}
	ATData, err := setteamInfo(filepath, awayteamName)
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
	HTData.updateBackgroundHometeamRatings(HomeErrFunc)
	ATData.updateBackgroundAwayteamRatings(AwayErrFunc)

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

// Search for a teams rating
//
// Error should be checked first before trying to use team struct
func Search(filepath string, teamName string) (Team, error) {
	teamInfo, err := checkRatings(filepath, teamName)
	if err != nil {
		return Team{}, err
	}
	tm := Team{}
	tm.Name = teamName
	tm.HomeRating, err = strconv.ParseFloat(teamInfo[1], 64)
	if err != nil {
		return Team{}, fmt.Errorf("%s: %v is not a number", teamName, tm.HomeRating)
	}
	tm.AwayRating, err = strconv.ParseFloat(teamInfo[2], 64)
	if err != nil {
		return Team{}, fmt.Errorf("%s: %v is not a number", teamName, tm.AwayRating)
	}
	tm.ContinuousPerformanceHome, err = strconv.Atoi(teamInfo[3])
	if err != nil {
		return Team{}, fmt.Errorf("%s: %v is not a number", teamName, tm.ContinuousPerformanceHome)
	}
	tm.ContinuousPerformanceAway, err = strconv.Atoi(teamInfo[4])
	if err != nil {
		return Team{}, fmt.Errorf("%s: %v is not a number", teamName, tm.ContinuousPerformanceAway)

	}

	return tm, nil
}
