package pi

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/arewabolu/csvmanager"
	"golang.org/x/exp/slices"
)

func TestNewRating(t *testing.T) {
	HT := &Team{
		Name:                      "Leicester",
		HomeRating:                1.6,
		AwayRating:                0.4,
		ContinuousPerformanceHome: 3,
	}
	HXG := expectedGoalIndividual(HT.HomeRating)
	goalScoredH := 4
	AT := &Team{
		Name:                      "Wolves",
		HomeRating:                0.3,
		AwayRating:                -1.2,
		ContinuousPerformanceAway: -1,
	}
	AXG := expectedGoalIndividual(AT.AwayRating)
	goalScoredA := 1
	errFunc := errorGDFunc(errorGD(goalDifference(goalScoredH, goalScoredA), expectedGoalDifference(HXG, AXG)))
	//t.Error(HT.ProvisionalRatingHome())
	HT.UpdateBackgroundHomeTeamRatings(errFunc)
	HT.UpdateContinuousPerformanceHome()
	AT.UpdateContinuousPerformanceAway()
	t.Error(AT.ProvisionalRatingAway())
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
		HxG := expectedGoalIndividual(HomeTeam.HomeRating)
		AxG := expectedGoalIndividual(AwayTeam.AwayRating)
		xGD := expectedGoalDifference(HxG, AxG)
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

		HomeTeam.UpdateBackgroundHomeTeamRatings(HomeErrFunc)
		AwayTeam.UpdateBackgroundAwayTeamRatings(AwayErrFunc)
		switch {
		case xGD >= 0 && GD > 0:
			HomeTeam.UpdateContinuousPerformanceHome()
			AwayTeam.UpdateContinuousPerformanceAway()
		case xGD > 0 && GD < 0:
			HomeTeam.ResetContinuousPerformanceHome()
			AwayTeam.ResetContinuousPerformanceAway()
		case xGD < 0 && GD > 0:
			AwayTeam.ResetContinuousPerformanceAway()
			HomeTeam.ResetContinuousPerformanceHome()
		case xGD <= 0 && GD < 0:
			AwayTeam.UpdateContinuousPerformanceAwayV2()
			HomeTeam.UpdateContinuousPerformanceHomeV2()
		case xGD > 0 || xGD < 0 && GD == 0:
			HomeTeam.ResetContinuousPerformanceHome()
			AwayTeam.ResetContinuousPerformanceAway()
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
	data, err := csvmanager.ReadCsv("./matchdata/fifa4x4Eng.csv", 0755, true)
	if err != nil {
		panic(err)
	}
	rows := data.Rows()
	for _, game := range rows {
		match := game.String()
		HomeTeam := &Team{Name: match[0]}
		AwayTeam := &Team{Name: match[1]}
		homeGoal, err := strconv.Atoi(match[2])
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		awayGoal, err := strconv.Atoi(match[3])
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		err = UpdateTeamRatings("./matchdata/ratingsfifa4x4Eng.csv", HomeTeam.Name, AwayTeam.Name, homeGoal, awayGoal)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	//UpdateTeamRatings("LIV", "NOR", 4, 1)
}

func TestSearchRankings(t *testing.T) {
	search := func(teamName, venue string) *Team {
		ratings, err := csvmanager.ReadCsv("./matchdata/ratingsfifa4x4Eng.csv", 0755, true)
		if err != nil {
			panic(err)
		}
		Team := &Team{}
		TeamCol := ratings.Col("TeamName").String()
		occurence := slices.Index(TeamCol, teamName)
		data := ratings.Row(occurence).String()
		Team.Name = teamName
		Team.HomeRating, err = strconv.ParseFloat(data[1], 64)
		if err != nil {
			panic(err)
		}
		Team.AwayRating, err = strconv.ParseFloat(data[2], 64)
		if err != nil {
			panic(err)
		}
		Team.ContinuousPerformanceHome, err = strconv.Atoi(data[3])
		if err != nil {
			panic(err)
		}
		Team.ContinuousPerformanceAway, err = strconv.Atoi(data[4])
		if err != nil {
			panic(err)
		}
		if venue == "away" {
			if Team.ContinuousPerformanceAway > 1 {
				Team.AwayRating = Team.ProvisionalRatingAwayV2()
			} else if Team.ContinuousPerformanceAway < -1 {
				Team.AwayRating = Team.ProvisionalRatingAway()
			}
		}
		if venue == "home" {
			if Team.ContinuousPerformanceHome > 1 {
				Team.HomeRating = Team.ProvisionalRatingHome()
			} else if Team.ContinuousPerformanceHome < -1 {
				Team.HomeRating = Team.ProvisionalRatingHomeV2()
			}
		}

		return Team
	}
	t.Error(search("BHA", "away"))
}

//
/*

switch {
	case HToccurence == -1 && AToccurence == -1:
		HTWrite := []string{HomeTeam.Name, fmt.Sprintf("%.2f", HomeTeam.HomeRating), fmt.Sprintf("%.2f", HomeTeam.AwayRating), fmt.Sprintf("%d", HomeTeam.ContinuousPerformanceHome), fmt.Sprintf("%d", HomeTeam.ContinuousPerformanceAway)}
		ATWrite := []string{AwayTeam.Name, fmt.Sprintf("%.2f", AwayTeam.HomeRating), fmt.Sprintf("%.2f", AwayTeam.AwayRating), fmt.Sprintf("%d", AwayTeam.ContinuousPerformanceHome), fmt.Sprintf("%d", AwayTeam.ContinuousPerformanceAway)}
		writeRatings := csvmanager.WriteFrame{
			Row:    true,
			Arrays: [][]string{HTWrite, ATWrite}, //
			File:   file,
		}
		writeRatings.WriteCSV()
		file.Close()
		t.Error("case 1")
	case HToccurence > -1 && AToccurence == -1:
		HTWrite := []string{HomeTeam.Name, fmt.Sprintf("%.2f", HomeTeam.HomeRating), fmt.Sprintf("%.2f", HomeTeam.AwayRating), fmt.Sprintf("%d", HomeTeam.ContinuousPerformanceHome), fmt.Sprintf("%d", HomeTeam.ContinuousPerformanceAway)}
		csvmanager.ReplaceRow("./matchdata/ratingsfifa4x4Eng.csv", 0755, HToccurence+1, HTWrite)
		ATWrite := []string{AwayTeam.Name, fmt.Sprintf("%.2f", AwayTeam.HomeRating), fmt.Sprintf("%.2f", AwayTeam.AwayRating), fmt.Sprintf("%d", AwayTeam.ContinuousPerformanceHome), fmt.Sprintf("%d", AwayTeam.ContinuousPerformanceAway)}
		writeRatings := csvmanager.WriteFrame{
			Row:    true,
			Arrays: [][]string{ATWrite}, //
			File:   file,
		}
		writeRatings.WriteCSV()
		file.Close()
		t.Error("case 2")
	case AToccurence > -1 && HToccurence == -1:
		HTWrite := []string{HomeTeam.Name, fmt.Sprintf("%.2f", HomeTeam.HomeRating), fmt.Sprintf("%.2f", HomeTeam.AwayRating), fmt.Sprintf("%d", HomeTeam.ContinuousPerformanceHome), fmt.Sprintf("%d", HomeTeam.ContinuousPerformanceAway)}
		writeRatings := csvmanager.WriteFrame{
			Row:    true,
			Arrays: [][]string{HTWrite}, //
			File:   file,
		}
		writeRatings.WriteCSV()
		file.Close()
		ATWrite := []string{AwayTeam.Name, fmt.Sprintf("%.2f", AwayTeam.HomeRating), fmt.Sprintf("%.2f", AwayTeam.AwayRating), fmt.Sprintf("%v", AwayTeam.ContinuousPerformanceHome), fmt.Sprintf("%v", AwayTeam.ContinuousPerformanceAway)}
		csvmanager.ReplaceRow("./matchdata/ratingsfifa4x4Eng.csv", 0755, AToccurence+1, ATWrite)
		t.Error("case 3")
	case HToccurence > -1 && AToccurence > -1:

		t.Error("case 4")
	}*/
