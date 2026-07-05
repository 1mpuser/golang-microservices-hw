package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/1mpuser/order/internal/app"
	"github.com/1mpuser/order/internal/config"
)

func main() {
	_ = godotenv.Load("order.env") //nolint:gosec // .env файл опционален — ошибка загрузки допустима

	config.MustLoad()

	if err := app.New(context.Background()).Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
