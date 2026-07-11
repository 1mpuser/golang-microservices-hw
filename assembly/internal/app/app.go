package app

import (
	"context"
	"log/slog"
)

// App управляет жизненным циклом AssemblyService (Kafka consumer + producer).
type App struct {
	di *diContainer
}

// New создаёт приложение и инициализирует зависимости.
func New(_ context.Context) *App {
	return &App{
		di: newDiContainer(),
	}
}

// Run запускает consumer OrderPaid и ждёт сигнала ОС для graceful shutdown.
//
// TODO(неделя 5, часть 4): реализовать —
//  1. logger.Init(config.AppConfig().Logger.Level).
//  2. Получить consumer из DI, запустить consumer.Run(ctx) в горутине.
//  3. Ждать SIGINT/SIGTERM (signal.NotifyContext), затем closer.CloseAll.
func (a *App) Run() error {
	slog.Info("запуск AssemblyService — TODO: реализовать часть 4")
	panic("не реализовано: часть 4 (запуск AssemblyService)")
}
