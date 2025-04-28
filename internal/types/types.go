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
	To            string   `json:"to" validate:"required"`
	Subject       string   `json:"subject"`
	Greeting      string   `json:"greeting"`
	Message       string   `json:"message" validate:"required"`
	ActionDetails string   `json:"actionDetails"`
	NextSteps     []string `json:"steps"`
	RecipientName string   `json:"recipientName"`
	SenderName    string   `json:"senderName"`
	SenderTitle   string   `json:"senderTitle"`
}
