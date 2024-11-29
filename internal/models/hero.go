package models

import (
	"errors"
	"fmt"
	utils2 "github.com/main_projects/bbot/epic_battle/internal/utils"
	"github.com/main_projects/bbot/epic_battle/pkg/config"
)

const (
	// Типи атак
	NormalAttack AttackType = iota
	WhirlwindAttack
	PreciseStrike
	MightyBlow
)

type AttackType int

type Attack struct {
	Type        AttackType
	HeadsCut    int
	Description string
}

type Hero struct {
	Name    string
	config  config.GameConfig
	Stunned bool
}

func NewHero(cfg config.GameConfig) (*Hero, error) {
	if err := validateHeroConfig(cfg); err != nil {
		return nil, fmt.Errorf("невалідна конфігурація героя: %w", err)
	}

	return &Hero{
		Name:    utils2.GenerateHeroName(),
		config:  cfg,
		Stunned: false,
	}, nil
}

func validateHeroConfig(cfg config.GameConfig) error {
	if cfg.MinHeroAttack <= 0 {
		return errors.New("мінімальна атака має бути більше 0")
	}
	if cfg.MaxHeroAttack <= cfg.MinHeroAttack {
		return errors.New("максимальна атака має бути більше мінімальної")
	}
	if cfg.SpecialAttackChance < 0 || cfg.SpecialAttackChance > 100 {
		return errors.New("шанс спеціальної атаки має бути між 0 та 100")
	}
	return nil
}

func (h *Hero) Attack() Attack {
	if h.Stunned {
		h.Stunned = false
		return Attack{
			Type:        NormalAttack,
			HeadsCut:    0,
			Description: "богатир оглушений і пропускає хід!",
		}
	}

	if utils2.RandomInt(1, 100) <= h.config.SpecialAttackChance {
		return h.specialAttack()
	}

	headsCut := utils2.RandomInt(h.config.MinHeroAttack, h.config.MaxHeroAttack)
	return Attack{
		Type:        NormalAttack,
		HeadsCut:    headsCut,
		Description: "звичайна атака мечем",
	}
}

func (h *Hero) specialAttack() Attack {
	attackType := utils2.RandomInt(1, 3)
	switch attackType {
	case 1:
		baseAttack := utils2.RandomInt(h.config.MinHeroAttack, h.config.MaxHeroAttack)
		headsCut := baseAttack + utils2.RandomInt(0, baseAttack)
		return Attack{
			Type:        WhirlwindAttack,
			HeadsCut:    headsCut,
			Description: "ВИХОР КЛИНКА",
		}
	case 2:
		headsCut := h.config.MaxHeroAttack + utils2.RandomInt(1, 2)
		return Attack{
			Type:        PreciseStrike,
			HeadsCut:    headsCut,
			Description: "ТОЧНИЙ УДАР",
		}
	default:
		baseAttack := utils2.RandomInt(h.config.MinHeroAttack, h.config.MaxHeroAttack)
		headsCut := int(float64(baseAttack) * 1.5)
		return Attack{
			Type:        MightyBlow,
			HeadsCut:    headsCut,
			Description: "МОГУТНІЙ УДАР",
		}
	}
}
