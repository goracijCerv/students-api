package types

type Student struct {
	Id       int64  `json:"id"`
	Name     string `json:"name" validate:"required"`
	LastName string `json:"lastname" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Number   int    `json:"number" validate:"required"`
	Age      int    `json:"age" validate:"required"`
}

type EmailData struct {
	Subject       string
	Greeting      string
	Message       string
	ActionDetails string
	NextSteps     []string
	RecipientName string
	SenderName    string
	SenderTitle   string
}

type EmailReciber struct {
	To string `json:"to" validate:"required"`
}
