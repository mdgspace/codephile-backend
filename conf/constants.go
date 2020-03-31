package conf

import "github.com/mdg-iitr/Codephile/models/types"

const (
	CODECHEF   = "codechef"
	CODEFORCES = "codeforces"
	HACKERRANK = "hackerrank"
	SPOJ       = "spoj"
)

var (
	codechefSubMutex   = types.SubmissionMutex{Website: CODECHEF}
	codeforcesSubMutex = types.SubmissionMutex{Website: CODECHEF}
	hackerrankSubMutex = types.SubmissionMutex{Website: CODECHEF}
	spojSubMutex       = types.SubmissionMutex{Website: CODECHEF}
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
func GetMutexForWebsiteSubmission(site string) *types.SubmissionMutex {
	switch site {
	case CODECHEF:
		return &codechefSubMutex
	case CODEFORCES:
		return &codeforcesSubMutex
	case HACKERRANK:
		return &hackerrankSubMutex
	case SPOJ:
		return &spojSubMutex
	default:
		return nil
	}
}
