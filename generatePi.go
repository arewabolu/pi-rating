package pi

//BuildPiforHometeam and BuildPiforAwayteam should only be used when generating ratings in real time.
//
//Otherwise you should use a file to manage the ratings provided by the makefile, Search and WritePi functions

// create default team structure
func Newteam(teamname string) Team {
	return Team{Name: teamname}
}

// generate a teams home ratings
func BuildPiforHometeam(hometeam Team, awayteam string, homeGoalScored, awayGoalScored int) Team {
	HxG := ExpectedGoalIndividual(hometeam.HomeRating)
	Awayteam := Newteam(awayteam)
	AxG := ExpectedGoalIndividual(Awayteam.HomeRating)
	xGD := ExpectedGoalDifference(HxG, AxG)
	GD := goalDifference(homeGoalScored, awayGoalScored)

	var HomeErrFunc float64
	errFunc := errorGDFunc(errorGD(GD, xGD))
	if xGD > float64(GD) {
		HomeErrFunc = -errFunc

	} else {
		HomeErrFunc = errFunc

	}
	hometeam.updateBackgroundHometeamRatings(HomeErrFunc)

	switch {
	case xGD >= 0 && GD > 0:
		hometeam.updateContinuousPerformanceHome()

	case xGD > 0 && GD < 0:
		hometeam.resetContinuousPerformanceHome()
		hometeam.updateContinuousPerformanceHomeV2()

	case xGD < 0 && GD > 0:
		hometeam.resetContinuousPerformanceHome()
		hometeam.updateContinuousPerformanceHome()
	case xGD <= 0 && GD < 0:

		hometeam.updateContinuousPerformanceHomeV2()
	case xGD > 0 || xGD < 0 && GD == 0:
		hometeam.resetContinuousPerformanceHome()
	}

	return hometeam
}

// generate a teams away ratings
func BuildPiforAwayteam(awayteam Team, hometeamName string, homeGoalScored, awayGoalScored int) Team {
	hometeam := Newteam(hometeamName)
	HxG := ExpectedGoalIndividual(hometeam.HomeRating)
	AxG := ExpectedGoalIndividual(awayteam.HomeRating)
	xGD := ExpectedGoalDifference(HxG, AxG)
	GD := goalDifference(homeGoalScored, awayGoalScored)

	var AwayErrFunc float64
	errFunc := errorGDFunc(errorGD(GD, xGD))
	if xGD > float64(GD) {
		AwayErrFunc = errFunc
	} else {
		AwayErrFunc = -errFunc
	}
	awayteam.updateBackgroundAwayteamRatings(AwayErrFunc)

	switch {
	case xGD >= 0 && GD > 0:
		awayteam.updateContinuousPerformanceAway()
	case xGD > 0 && GD < 0:
		awayteam.resetContinuousPerformanceAway()
		awayteam.updateContinuousPerformanceAwayV2()
	case xGD < 0 && GD > 0:
		awayteam.resetContinuousPerformanceAway()
		awayteam.updateContinuousPerformanceAway()
	case xGD <= 0 && GD < 0:
		awayteam.updateContinuousPerformanceAwayV2()
	case xGD > 0 || xGD < 0 && GD == 0:
		awayteam.resetContinuousPerformanceAway()
	}

	return awayteam
}
