package internal

// User represents a farmer
type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Uploads  []string `json:"uploads"`
	Calendar []string `json:"calendar"`
}