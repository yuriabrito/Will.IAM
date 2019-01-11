package models

// AM represents an item from /am []
type AM struct {
	Prefix   string `json:"prefix"`
	Alias    string `json:"alias"`
	Owner    bool   `json:"owner"`
	Lender   bool   `json:"lender"`
	Complete bool   `json:"complete"`
}
