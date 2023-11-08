package pi

import (
	"strconv"
	"testing"

	"github.com/arewabolu/csvmanager"
)

func TestUpdateTeamRatings(t *testing.T) {
	data, err := csvmanager.ReadCsv("./matchdata/fifa4x4Eng.csv", true)
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
		_, _, err = UpdateTeamRatings("./matchdata/ratingsfifa4x4Eng.csv", HomeTeam.Name, AwayTeam.Name, homeGoal, awayGoal)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	//UpdateTeamRatings("LIV", "NOR", 4, 1)
}

func TestSearchRankings(t *testing.T) {
	hT, _ := Search("./matchdata/ratingsfifa4x4Eng.csv", "LEI")
	aT, _ := Search("./matchdata/ratingsfifa4x4Eng.csv", "NU")
	t.Error(hT.HomeRating, aT.AwayRating)
}

func TestBuilderv2(t *testing.T) {
	data, err := csvmanager.ReadCsv("./matchdata/fifa4x4Eng.csv", true)
	if err != nil {
		panic(err)
	}
	ht := Newteam("INDIANA")
	at := Newteam("OPP")

	rows := data.Rows()
	for _, game := range rows {
		match := game.String()
		if match[0] == "INDIANA" {
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

			ht = BuildPiforHometeamV2(ht, &at, homeGoal, awayGoal)
		}
	}
	t.Error(ht)
}
