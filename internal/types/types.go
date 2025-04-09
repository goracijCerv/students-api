package types

type Student struct {
	Id       int
	Name     string `validate:"required"`
	LastName string `validate:"required"`
	Email    string `validate:"required"`
	Number   int    `validate:"required"`
	Age      int    `validate:"required"`
}
