package conf

import (
	"fmt"
	"strings"
)

const (
	CODECHEF   = "codechef"
	CODEFORCES = "codeforces"
	HACKERRANK = "hackerrank"
	SPOJ       = "spoj"
	LEETCODE   = "leetcode"
)

var ValidSites = []string{HACKERRANK, CODECHEF, CODEFORCES, SPOJ, LEETCODE}

func GetRegexSite(site string) string {
	switch site {
	case CODECHEF:
		return "https://www.codechef.com"
	case CODEFORCES:
		return "http://codeforces.com"
	case HACKERRANK:
		return "https://www.hackerrank.com"
	case SPOJ:
		return "https://www.spoj.com"
	case LEETCODE:
		return "https://leetcode.com/"
	}
	return " "
}

func IsSiteValid(s string) bool {
	for _, vs := range ValidSites {
		if s == vs {
			return true
		}
	}
	return false
}

func GetSiteFromURL(url string) (string, error) {
	for _, vs := range ValidSites {
		if strings.Contains(url, vs) {
			return vs, nil
		}
	}
	return url, fmt.Errorf("unrecognised platform URL: %s", url)
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

var ValidPaths = []string{"username", "fullname", "handle.codechef", "handle.codeforces", "handle.hackerearth", "handle.hackerrank", "handle.spoj"}

func IsPathValid(s string) bool {
	for _, path := range ValidPaths {
		if path == s {
			return true
		}
	}
	return false
}
