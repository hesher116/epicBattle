package main

import (
	"context"
	"fmt"
	"github.com/main_projects/bbot/epic_battle/internal/battle"
	"github.com/main_projects/bbot/epic_battle/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Створюємо контекст з можливістю скасування
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Канал для результату битви
	done := make(chan error, 1)

	// Запускаємо битву в горутині
	go func() {
		logger.Info("Починається епічна битва!")
		done <- battle.StartEpicBattle(ctx)
	}()

	// Очікуємо або завершення битви, або сигналу завершення
	select {
	case <-shutdown:
		logger.Warn("Отримано сигнал завершення")
		cancel()
		// Даємо час на graceful shutdown
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			logger.Error("Таймаут завершення битви")
		}
	case err := <-done:
		if err != nil {
			logger.Error(fmt.Sprintf("Помилка: %v", err))
			os.Exit(1)
		}
	}

	logger.Info("Епічна битва завершена!")
}
