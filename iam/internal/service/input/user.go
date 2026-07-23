package input

// RegisterInput — валидированный вход метода Register (пришёл из API-слоя).
type RegisterInput struct {
	Login    string
	Password string
}

// LoginInput — валидированный вход метода Login (пришёл из API-слоя).
type LoginInput struct {
	Login    string
	Password string
}
