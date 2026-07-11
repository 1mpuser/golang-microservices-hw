package main

import (
	"context"
	"log/slog"

	"github.com/1mpuser/assembly/internal/app"
	"github.com/1mpuser/assembly/internal/config"
)

func main() {
	config.MustLoad()

	if err := app.New(context.Background()).Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
