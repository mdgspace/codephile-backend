package types

type ActivityGraph []SubmissionCount
type SubmissionCount struct {
	Correct   int    `json:"correct"`
	Total     int    `json:"total"`
	CreatedAt string `bson:"_id" json:"created_at"`
}

type StatusCounts struct {
	StatusCorrect             int `bson:"ac_count"`
	StatusWrongAnswer         int `bson:"wa_count"`
	StatusCompilationError    int `bson:"ce_count"`
	StatusRuntimeError        int `bson:"re_count"`
	StatusTimeLimitExceeded   int `bson:"tle_count"`
	StatusMemoryLimitExceeded int `bson:"mle_count"`
	StatusPartial             int `bson:"ptl_count"`
}
