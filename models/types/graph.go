package types

type ActivityGraph []SubmissionCount
type SubmissionCount struct {
	Correct   int    `json:"correct"`
	Total     int    `json:"total"`
	CreatedAt string `bson:"_id" json:"created_at"`
}

type StatusCounts struct {
	StatusCorrect             int `bson:"ac_count" json:"ac"`
	StatusWrongAnswer         int `bson:"wa_count" json:"wa"`
	StatusCompilationError    int `bson:"ce_count" json:"ce"`
	StatusRuntimeError        int `bson:"re_count" json:"re"`
	StatusTimeLimitExceeded   int `bson:"tle_count" json:"tle"`
	StatusMemoryLimitExceeded int `bson:"mle_count" json:"mle"`
	StatusPartial             int `bson:"ptl_count" json:"ptl"`
}
