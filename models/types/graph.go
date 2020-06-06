package types

type ActivityGraph []SubmissionCount
type SubmissionCount struct {
	Correct   int    `json:"correct"`
	Total     int    `json:"total"`
	CreatedAt string `bson:"_id" json:"created_at"`
}
