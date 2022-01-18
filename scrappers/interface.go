package scrappers

import (
	"context"
	"errors"
	"time"

	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scrappers/codechef"
	"github.com/mdg-iitr/Codephile/scrappers/codeforces"
	"github.com/mdg-iitr/Codephile/scrappers/hackerrank"
	"github.com/mdg-iitr/Codephile/scrappers/spoj"
)

type Scrapper interface {
	CheckHandle() (bool, error)
	GetSubmissions(after time.Time) []types.Submission
	GetProfileInfo() types.ProfileInfo
}

func NewScrapper(site string, handle string, ctx context.Context) (Scrapper, error) {
	if handle == "" {
		return nil, HandleNotFoundError
	}
	switch site {
	case CODEFORCES:
		return codeforces.Scrapper{Handle: handle, Context: ctx}, nil
	case CODECHEF:
		return codechef.Scrapper{Handle: handle, Context: ctx}, nil
	case HACKERRANK:
		return hackerrank.Scrapper{Handle: handle, Context: ctx}, nil
	case SPOJ:
		return spoj.Scrapper{Handle: handle, Context: ctx}, nil
	case LEETCODE:
		return spoj.Scrapper{Handle: handle, Context: ctx}, nil
	default:
		return nil, errors.New("site invalid")
	}
}
