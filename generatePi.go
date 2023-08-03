package pi

//BuildPiforHometeam and BuildPiforAwayteam should only be used when generating ratings in real time.
//
//Otherwise you should use a file to manage the ratings provided by the makefile, Search and WritePi functions

// create default team structure
func Newteam(teamname string) Team {
	return Team{Name: teamname}
}

// generate a teams home ratings, assuming every opposition is new and has no previous rating
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

// generate a teams away ratings, assuming every opposition is new and has no previous rating
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

// This version caters for the opposition teams rating and how likely any team will score
// against the given team.
func BuildPiforHometeamV2(hometeam Team, awayteam *Team, homeGoalScored, awayGoalScored int) Team {
	HxG := ExpectedGoalIndividual(hometeam.HomeRating)
	//	Awayteam := Newteam(awayteam)
	AxG := ExpectedGoalIndividual(awayteam.AwayRating)
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

	var awayErrFunc float64
	errFunc2 := errorGDFunc(errorGD(GD, xGD))
	if xGD > float64(GD) {
		awayErrFunc = errFunc2
	} else {
		awayErrFunc = -errFunc2
	}
	awayteam.updateBackgroundAwayteamRatings(awayErrFunc)

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

// This version caters for the opposition teams rating and how likely any team will score
// against the given team.
func BuildPiforAwayteamv2(awayteam Team, hometeam *Team, homeGoalScored, awayGoalScored int) Team {
	//hometeam := Newteam(hometeamName)
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

	var HomeErrFunc float64
	if xGD > float64(GD) {
		HomeErrFunc = -errFunc
	} else {
		HomeErrFunc = errFunc
	}
	hometeam.updateBackgroundHometeamRatings(HomeErrFunc)

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
