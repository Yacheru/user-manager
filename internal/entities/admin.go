package entities

type Task struct {
	ID          *int    `json:"-" db:"id" swaggerignore:"true"`
	TaskID      *string `json:"uuid" db:"uuid" swaggerignore:"true"`
	Title       string  `json:"title" db:"title"`
	Description string  `json:"description" db:"description"`
	Reward      int     `json:"reward" db:"reward"`
}

type Code struct {
	Code   string `json:"code" bson:"code" binding:"required"`
	Reward int    `json:"reward" bson:"reward" binding:"required"`
}
