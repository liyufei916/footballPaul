package utils

const (
	PointsExactScore              = 10
	PointsCorrectResultAndOneScore = 5 // 猜对胜平负且其中一方进球数正确
	PointsCorrectResult           = 3 // 只猜对胜平负
	PointsIncorrect               = 0
)

func CalculatePoints(predictedHome, predictedAway, actualHome, actualAway int) int {
	// 规则1：比分完全准确，得10分
	if predictedHome == actualHome && predictedAway == actualAway {
		return PointsExactScore
	}

	predictedResult := getMatchResult(predictedHome, predictedAway)
	actualResult := getMatchResult(actualHome, actualAway)

	if predictedResult != actualResult {
		// 规则4：胜平负猜错，得0分
		return PointsIncorrect
	}

	// 规则2：猜对胜平负，且其中一方进球数正确，得5分
	if predictedHome == actualHome || predictedAway == actualAway {
		return PointsCorrectResultAndOneScore
	}

	// 规则3：只猜对胜平负，但没猜对任何一方进球数，得3分
	return PointsCorrectResult
}

func getMatchResult(homeScore, awayScore int) string {
	if homeScore > awayScore {
		return "home_win"
	} else if homeScore < awayScore {
		return "away_win"
	}
	return "draw"
}
