package models

type Users struct {
	Model
	Username string
	Password string
	Role     string
}
