package user

import "time"

type User struct {
	ID        string
	Email     string
	Name      string
	Role      string
	Status    string
	Photo     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
