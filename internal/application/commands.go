package application

type CreateUserCommand struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type CreateUserResult struct {
	Token string `json:"token"`
}
