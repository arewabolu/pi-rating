package pi

import "testing"

// A run of good results must raise a provisional rating and a run of poor
// results must lower it, at both venues.
func TestProvisionalRatingDirection(t *testing.T) {
	const base = 0.5
	cases := []struct {
		name  string
		team  Team
		venue string
		up    bool
	}{
		{"good home run", Team{HomeRating: base, ContinuousPerformanceHome: 4}, "home", true},
		{"poor home run", Team{HomeRating: base, ContinuousPerformanceHome: -4}, "home", false},
		{"good away run", Team{AwayRating: base, ContinuousPerformanceAway: -4}, "away", true},
		{"poor away run", Team{AwayRating: base, ContinuousPerformanceAway: 4}, "away", false},
	}
	for _, c := range cases {
		p := c.team.ProvisionalRating(c.venue)
		got := p.HomeRating
		if c.venue == "away" {
			got = p.AwayRating
		}
		if (got > base) != c.up || got == base {
			t.Errorf("%s: provisional %v from %v, want it to go %s", c.name, got, base, map[bool]string{true: "up", false: "down"}[c.up])
		}
	}
}
