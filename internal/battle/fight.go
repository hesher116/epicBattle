package battle

import (
	"context"
	"fmt"
	"github.com/main_projects/bbot/epic_battle/internal/models"
	"github.com/main_projects/bbot/epic_battle/pkg/config"
	"github.com/main_projects/bbot/epic_battle/pkg/errors"
	"github.com/main_projects/bbot/epic_battle/pkg/logger"
)

func StartEpicBattle(ctx context.Context) error {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return fmt.Errorf("помилка конфігурації: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("невалідна конфігурація: %w", err)
	}

	hero, err := models.NewHero(cfg)
	if err != nil {
		return fmt.Errorf("помилка створення героя: %w", err)
	}

	dragon, err := models.NewDragon(cfg)
	if err != nil {
		return fmt.Errorf("помилка створення дракона: %w", err)
	}

	battle, err := models.NewBattle(hero, dragon, cfg)
	if err != nil {
		return fmt.Errorf("помилка створення битви: %w", err)
	}

	logger.Info("Битва починається", logger.Fields{
		"hero_name":    hero.Name,
		"dragon_heads": dragon.Heads,
		"config":       cfg,
	})

	for {
		select {
		case <-ctx.Done():
			logger.Info("Битва перервана", logger.Fields{
				"reason": "context cancelled",
				"round":  battle.Round,
			})
			return ctx.Err()
		default:
			err := battle.ExecuteRound()
			if err != nil {
				switch err.(type) {
				case *errors.HeroVictoryError:
					logger.Info("Богатир переміг!", logger.Fields{
						"rounds":       battle.Round,
						"dragon_heads": dragon.Heads,
					})
					fmt.Println("🎉 УРА! Богатир переміг! Царство врятоване!")
				case *errors.DragonVictoryError:
					logger.Error("Змій переміг!", logger.Fields{
						"rounds":       battle.Round,
						"dragon_heads": dragon.Heads,
					})
					fmt.Println("😭 О ні! Змій переміг! Богатир втікає!")
				case *errors.DragonPhobiaError:
					logger.Warn("Змій втік через фобію!", logger.Fields{
						"rounds":       battle.Round,
						"dragon_heads": dragon.Heads,
					})
					fmt.Println("😅 Ого! Змій раптово втік. Перемога?")
				default:
					logger.Error("Несподівана помилка", logger.Fields{
						"error":        err,
						"rounds":       battle.Round,
						"dragon_heads": dragon.Heads,
					})
					fmt.Println("🤔 Щось пішло не так... Битва закінчилась незрозуміло.")
				}
				return err
			}
		}
	}
}
