package pi

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/arewabolu/csvmanager"
	"golang.org/x/exp/slices"
)

// Checks if a teams is listed in the rating table
// and returns the team information if it exists.

func Addteam(filepath, teamName string) error {
	teamName = strings.ToUpper(strings.TrimSpace(teamName))
	if teamName == "" {
		return errors.New("please state the name of the team")
	}
	_, err := os.Stat(filepath)
	if errors.Is(err, os.ErrNotExist) {
		Makefile(filepath)
	}
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	wr := csv.NewWriter(file)
	defer wr.Flush()

	err = wr.Write([]string{teamName, "0", "0", "0", "0"})
	if err != nil {
		return err
	}
	return nil
}

func checkRatings(filepath string, teamName string) ([]string, error) {
	ratings, err := csvmanager.ReadCsv(filepath, 0755, true)
	if err != nil {
		return nil, fmt.Errorf("%s not found", filepath)
	}

	if !slices.Contains(ratings.ListHeaders(), "teamName") {
		return nil, err
	}
	teamCol := ratings.Col("teamName").String()
	index := slices.Index(teamCol, teamName)
	if index != -1 {
		occurence := ratings.Row(index).String()
		occurence = append(occurence, fmt.Sprint(index))
		return occurence, nil
	}
	return nil, fmt.Errorf("%s not registered. please make sure team is already registered", teamName)
}

func Makefile(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	f := csvmanager.WriteFrame{
		Headers: []string{"teamName", "HomeRating", "AwayRating", "ContinuousHomePerformance", "ContinuousAwayPerformance"},
		File:    file,
	}
	err = f.WriteCSV()
	if err != nil {
		return err
	}
	return nil
}

func setteamInfo(filepath string, teamName string) (team, error) {
	teamData, err := checkRatings(filepath, teamName)
	if err != nil {
		return team{}, err
	}
	team := team{Name: teamName}
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

// Write updated rating to rating file
func (t *team) WriteRatings(filepath string) error {
	teamData := []string{t.Name, fmt.Sprintf("%.2f", t.HomeRating), fmt.Sprintf("%.2f", t.AwayRating), fmt.Sprintf("%d", t.ContinuousPerformanceHome), fmt.Sprintf("%d", t.ContinuousPerformanceAway)}
	err := csvmanager.ReplaceRow(filepath, 0755, t.index+1, teamData)
	return err
}
