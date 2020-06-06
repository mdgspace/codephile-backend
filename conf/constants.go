package conf

const (
	CODECHEF   = "codechef"
	CODEFORCES = "codeforces"
	HACKERRANK = "hackerrank"
	SPOJ       = "spoj"
)

var ValidSites = []string{HACKERRANK, CODECHEF, CODEFORCES, SPOJ}

func IsSiteValid(s string) bool {
	for _, vs := range ValidSites {
		if s == vs {
			return true
		}
	}
	return false
}

const (
	StatusCorrect             = "AC"
	StatusWrongAnswer         = "WA"
	StatusCompilationError    = "CE"
	StatusRuntimeError        = "RE"
	StatusTimeLimitExceeded   = "TLE"
	StatusMemoryLimitExceeded = "MLE"
	StatusPartial             = "PTL"
)
