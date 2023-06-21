package pi

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/arewabolu/csvmanager"
	"golang.org/x/exp/slices"
)

const ratingFilepath = "./matchdata/ratingsfifa4x4Eng.csv"

func TestNewRating(t *testing.T) {
	HT := &Team{
		Name:                      "Leicester",
		HomeRating:                1.6,
		AwayRating:                0.4,
		ContinuousPerformanceHome: 3,
	}
	HXG := ExpectedGoalIndividual(HT.HomeRating)
	goalScoredH := 4
	AT := &Team{
		Name:                      "Wolves",
		HomeRating:                0.3,
		AwayRating:                -1.2,
		ContinuousPerformanceAway: -1,
	}
	AXG := ExpectedGoalIndividual(AT.AwayRating)
	goalScoredA := 1
	errFunc := errorGDFunc(errorGD(goalDifference(goalScoredH, goalScoredA), ExpectedGoalDifference(HXG, AXG)))
	//t.Error(HT.ProvisionalRatingHome())
	HT.updateBackgroundHomeTeamRatings(errFunc)
	HT.updateContinuousPerformanceHome()
	AT.updateContinuousPerformanceAway()
	t.Error(AT.provisionalRatingAway())
}

func TestReader(t *testing.T) {
	data, err := csvmanager.ReadCsv("./matchdata/fifa4x4Eng.csv", 0755, true)
	if err != nil {
		panic(err)
	}
	rows := data.Rows()
	for i := 0; i < len(rows); i++ {

		ratings, err := csvmanager.ReadCsv("./matchdata/ratingsfifa4x4Eng.csv", 0755, true)
		if err != nil {
			panic(err)
		}
		TeamCol := ratings.Col("TeamName").String()
		game := rows[i].String()
		homeGoal, err := strconv.Atoi(game[2])
		if err != nil {
			return
		}
		awayGoal, _ := strconv.Atoi(game[3])

		HomeTeam := &Team{Name: game[0]}
		AwayTeam := &Team{Name: game[1]}
		HToccurence := slices.Index(TeamCol, HomeTeam.Name)
		AToccurence := slices.Index(TeamCol, AwayTeam.Name)

		if HToccurence > -1 {
			HTData := ratings.Row(HToccurence + 1).String()
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
			HTData := ratings.Row(AToccurence).String()
			AwayTeam.HomeRating, err = strconv.ParseFloat(HTData[1], 64)
			if err != nil {
				panic(err)
			}
			AwayTeam.AwayRating, err = strconv.ParseFloat(HTData[2], 64)
			if err != nil {
				panic(err)
			}
			AwayTeam.ContinuousPerformanceHome, err = strconv.Atoi(HTData[3])
			if err != nil {
				panic(err)
			}
			AwayTeam.ContinuousPerformanceAway, err = strconv.Atoi(HTData[4])
			if err != nil {
				panic(err)
			}
		}
		HxG := ExpectedGoalIndividual(HomeTeam.HomeRating)
		AxG := ExpectedGoalIndividual(AwayTeam.AwayRating)
		xGD := ExpectedGoalDifference(HxG, AxG)
		GD := goalDifference(homeGoal, awayGoal)
		var HomeErrFunc, AwayErrFunc float64
		errFunc := errorGDFunc(errorGD(GD, xGD))
		if xGD > float64(GD) {
			HomeErrFunc = -errFunc
			AwayErrFunc = errFunc

		} else {
			HomeErrFunc = errFunc
			AwayErrFunc = -errFunc
		}
		if HomeTeam.Name == "LIV" {
			t.Error("before @Home", HomeTeam, ratings.Row(HToccurence).String())
		}
		if AwayTeam.Name == "LIV" {
			t.Error("before away", AwayTeam, ratings.Row(AToccurence).String())
		}

		HomeTeam.updateBackgroundHomeTeamRatings(HomeErrFunc)
		AwayTeam.updateBackgroundAwayTeamRatings(AwayErrFunc)
		switch {
		case xGD >= 0 && GD > 0:
			HomeTeam.updateContinuousPerformanceHome()
			AwayTeam.updateContinuousPerformanceAway()
		case xGD > 0 && GD < 0:
			HomeTeam.resetContinuousPerformanceHome()
			AwayTeam.resetContinuousPerformanceAway()
		case xGD < 0 && GD > 0:
			AwayTeam.resetContinuousPerformanceAway()
			HomeTeam.resetContinuousPerformanceHome()
		case xGD <= 0 && GD < 0:
			AwayTeam.updateContinuousPerformanceAwayV2()
			HomeTeam.updateContinuousPerformanceHomeV2()
		case xGD > 0 || xGD < 0 && GD == 0:
			HomeTeam.resetContinuousPerformanceHome()
			AwayTeam.resetContinuousPerformanceAway()
		}

		//file, _ := os.OpenFile("./matchdata/ratingsfifa4x4Eng.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
		HTWrite := []string{HomeTeam.Name, fmt.Sprintf("%.2f", HomeTeam.HomeRating), fmt.Sprintf("%.2f", HomeTeam.AwayRating), fmt.Sprintf("%d", HomeTeam.ContinuousPerformanceHome), fmt.Sprintf("%d", HomeTeam.ContinuousPerformanceAway)}
		ATWrite := []string{AwayTeam.Name, fmt.Sprintf("%.2f", AwayTeam.HomeRating), fmt.Sprintf("%.2f", AwayTeam.AwayRating), fmt.Sprintf("%d", AwayTeam.ContinuousPerformanceHome), fmt.Sprintf("%d", AwayTeam.ContinuousPerformanceAway)}
		if HomeTeam.Name == "LIV" {
			t.Error("@Home", HTWrite, xGD, GD, ratings.Row(0).String())
		}
		if AwayTeam.Name == "LIV" {
			t.Error("away", ATWrite, xGD, GD, AToccurence)
		}
		csvmanager.ReplaceRow("./matchdata/ratingsfifa4x4Eng.csv", 0755, HToccurence+1, HTWrite)
		csvmanager.ReplaceRow("./matchdata/ratingsfifa4x4Eng.csv", 0755, AToccurence+1, ATWrite)

	}

}

func TestUpdateTeamRatings(t *testing.T) {
	_, AT, _ := UpdateTeamRatings("./matchdata/ratingsfifa4x4Eng.csv", "EVE", "AVL", 2, 2)
	t.Error(AT)
	AT.WriteRatings("./matchdata/ratingsfifa4x4Eng.csv")
}

func TestSearchRankings(t *testing.T) {
	hT, err := Search("./matchdata/ratingsfifa4x4Eng.csv", "EVE")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	nwHT := hT.ProvisionalRating("home")
	aT, _ := Search("./matchdata/ratingsfifa4x4Eng.csv", "MCI")
	t.Error(nwHT.HomeRating, aT.AwayRating)
}

func TestCheckRatings(t *testing.T) {
	team, _ := CheckRatings(ratingFilepath, "EVE")
	t.Error(team)
}

func TestSetInfo(t *testing.T) {
	team, _ := SetTeamInfo(ratingFilepath, "MCI")
	t.Error(team)
}
