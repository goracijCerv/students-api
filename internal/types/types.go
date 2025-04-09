package types

type Student struct {
	Id       int64
	Name     string `validate:"required"`
	LastName string `validate:"required"`
	Email    string `validate:"required"`
	Number   int    `validate:"required"`
	Age      int    `validate:"required"`
}
