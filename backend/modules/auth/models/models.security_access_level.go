package models

// This constant will be moved to the service layer function that assigns
// SecurityAccessLevel during account creation. Leave here to not forget
const (
	Official  string = "official"
	Secret    string = "secret"
	TopSecret string = "top_secret"
)

// SecurityAccessPolicy models the levels of privileges used by accounts and documents.
// These privilege are used for classifying accounts who have access to them.
type SecurityAccessLevel struct {
	ClassificationLevel string `db:"classification_level" json:"classification_level"`
}

// SecurityAccessLevel will be it's own table in the DB
