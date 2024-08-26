package dbmodel

import "time"

type Note struct {
	Id        int       `db:"id"`
	Username  string    `db:"username"`
	Title     string    `db:"title"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}
