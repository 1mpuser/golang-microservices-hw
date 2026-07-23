package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/1mpuser/iam/internal/app"
	"github.com/1mpuser/iam/internal/config"
)

func main() {
	_ = godotenv.Load("iam.env") //nolint:gosec // .env-файл опционален

	config.MustLoad()

	if err := app.New(context.Background()).Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
