package main

import (
	"context"
	"log/slog"

	"github.com/1mpuser/payment/internal/app"
	"github.com/1mpuser/payment/internal/config"
)

func main() {
	config.MustLoad()

	if err := app.New(context.Background()).Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
