package pi

import (
	"strconv"
	"testing"

	"github.com/arewabolu/csvmanager"
)

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
