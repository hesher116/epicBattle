package models

import (
	"errors"
	"fmt"
	"github.com/main_projects/bbot/epic_battle/internal/utils"
	"github.com/main_projects/bbot/epic_battle/pkg/config"
)

const (
	FirebreathMultiplier = 2
	MaxChancePercentage  = 100
	MinHeads             = 1
)

type DragonAbility struct {
	Used        bool
	Description string
	Effect      string
}

type Dragon struct {
	Heads      int
	config     config.GameConfig
	Firebreath DragonAbility
	Stun       DragonAbility
}

func NewDragon(cfg config.GameConfig) (*Dragon, error) {
	if err := validateDragonConfig(cfg); err != nil {
		return nil, fmt.Errorf("невалідна конфігурація дракона: %w", err)
	}

	return &Dragon{
		Heads:      utils.RandomInt(cfg.InitialDragonHeadsMin, cfg.InitialDragonHeadsMax),
		config:     cfg,
		Firebreath: DragonAbility{Used: false},
		Stun:       DragonAbility{Used: false},
	}, nil
}

func validateDragonConfig(cfg config.GameConfig) error {
	if cfg.InitialDragonHeadsMin < MinHeads {
		return fmt.Errorf("мінімальна кількість голів має бути не менше %d", MinHeads)
	}
	if cfg.FirebreathChance < 0 || cfg.FirebreathChance > MaxChancePercentage {
		return errors.New("шанс вогняного подиху має бути між 0 та 100")
	}
	if cfg.StunChance < 0 || cfg.StunChance > MaxChancePercentage {
		return errors.New("шанс оглушення має бути між 0 та 100")
	}
	return nil
}

func (d *Dragon) UseAbilities(hero *Hero) DragonAbility {
	// Перевіряємо шанс вогняного подиху
	if utils.RandomInt(1, 100) <= d.config.FirebreathChance {
		d.Firebreath.Used = true
		return DragonAbility{
			Used:        true,
			Description: "🔥 Змій використовує ВОГНЯНИЙ ПОДИХ",
			Effect:      "регенерація подвоєна!",
		}
	}

	// Перевіряємо шанс оглушення
	if utils.RandomInt(1, 100) <= d.config.StunChance {
		d.Stun.Used = true
		hero.Stunned = true
		return DragonAbility{
			Used:        true,
			Description: "💫 Змій використовує ОГЛУШЛИВИЙ РЕВ",
			Effect:      "богатир пропустить наступний хід!",
		}
	}

	return DragonAbility{Used: false}
}

func (d *Dragon) RegenerateHeads(cutHeads int, firebreathActive bool) int {
	multiplier := 1
	if firebreathActive {
		multiplier = FirebreathMultiplier
		d.Firebreath.Used = false
	}

	regeneratedHeads := 0
	for i := 0; i < cutHeads; i++ {
		chance := utils.RandomInt(1, 100)

		switch {
		case chance <= d.config.NoRegenChance:
			// нічого не відростає
		case chance <= d.config.OneHeadChance:
			regeneratedHeads += 1 * multiplier
		case chance <= d.config.TwoHeadChance:
			regeneratedHeads += 2 * multiplier
		default:
			regeneratedHeads += 3 * multiplier
		}
	}
	return regeneratedHeads
}

func (d *Dragon) ResetAbilities() {
	d.Firebreath.Used = false
	d.Stun.Used = false
}
