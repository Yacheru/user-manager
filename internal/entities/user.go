package entities

type User struct {
	ID       *int   `json:"-" db:"id"`
	UserID   string `json:"uuid" db:"uuid"`
	Name     string `json:"name" db:"name"`
	Points   int    `json:"points" db:"points"`
	Referral bool   `json:"referral" db:"referral"`
	Tasks    []Task `json:"tasks" db:"tasks"`
}

type Leaderboard struct {
	UserID   string `json:"uuid" db:"uuid"`
	Name     string `json:"name" db:"name"`
	Points   int    `json:"points" db:"points"`
	Referral bool   `json:"referral" db:"referral"`
}

type NewUser struct {
	UserID *string `json:"id" db:"uuid" swaggerignore:"true"`
	Name   string  `json:"name" db:"name"`
}
