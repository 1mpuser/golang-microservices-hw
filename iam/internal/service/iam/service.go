package iam

// service — сервисный слой IAM. Зависит только от абстракций из deps.go
// (UserRepository — Postgres, SessionRepository — Redis). Методы Register/Login/
// Whoami/Logout/GetUser реализованы в отдельных файлах этого пакета.
type service struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewService(userRepo UserRepository, sessionRepo SessionRepository) *service {
	return &service{
		userRepo,
		sessionRepo,
	}
}
