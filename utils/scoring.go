package utils

type MatchResult string

const (
	HomeWin MatchResult = "home_win"
	AwayWin MatchResult = "away_win"
	Draw    MatchResult = "draw"
)

const (
	PointsExactScore      = 10
	PointsGoalDifference  = 7
	PointsCorrectWinner   = 5
	PointsOneScoreCorrect = 3
	PointsIncorrect       = 0
)

func CalculatePoints(predictedHome, predictedAway, actualHome, actualAway int) int {
	if predictedHome == actualHome && predictedAway == actualAway {
		return PointsExactScore
	}

	predictedResult := getMatchResult(predictedHome, predictedAway)
	actualResult := getMatchResult(actualHome, actualAway)

	if predictedResult != actualResult {
		if predictedHome == actualHome || predictedAway == actualAway {
			return PointsOneScoreCorrect
		}
		return PointsIncorrect
	}

	predictedDiff := abs(predictedHome - predictedAway)
	actualDiff := abs(actualHome - actualAway)

	if predictedDiff == actualDiff {
		return PointsGoalDifference
	}

	return PointsCorrectWinner
}

func getMatchResult(homeScore, awayScore int) MatchResult {
	if homeScore > awayScore {
		return HomeWin
	} else if homeScore < awayScore {
		return AwayWin
	}
	return Draw
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
