package models

import (
	"fmt"
	"github.com/main_projects/bbot/epic_battle/pkg/config"
	"github.com/main_projects/bbot/epic_battle/pkg/errors"
	"github.com/main_projects/bbot/epic_battle/pkg/logger"
)

type Battle struct {
	Hero   *Hero
	Dragon *Dragon
	Round  int
	config config.GameConfig
}

func NewBattle(hero *Hero, dragon *Dragon, cfg config.GameConfig) (*Battle, error) {
	// Валідуємо вхідні параметри
	if hero == nil {
		return nil, fmt.Errorf("герой не може бути nil")
	}
	if dragon == nil {
		return nil, fmt.Errorf("дракон не може бути nil")
	}

	// Перевіряємо початкові умови
	if dragon.Heads <= 0 {
		return nil, fmt.Errorf("дракон не може мати %d голів", dragon.Heads)
	}
	if dragon.Heads >= cfg.MaxDragonHeads {
		return nil, fmt.Errorf("дракон не може мати більше %d голів", cfg.MaxDragonHeads)
	}

	return &Battle{
		Hero:   hero,
		Dragon: dragon,
		Round:  0,
		config: cfg,
	}, nil
}

func (b *Battle) ExecuteRound() error {
	logger.Info("Початок раунду", logger.Fields{
		"round":        b.Round,
		"dragon_heads": b.Dragon.Heads,
		"hero":         b.Hero.Name,
	})
	b.Round++
	fmt.Printf("\n🗡️ Раунд %d: %s готується до атаки...\n", b.Round, b.Hero.Name)

	// Змій використовує здібності
	ability := b.Dragon.UseAbilities(b.Hero)
	if ability.Used {
		fmt.Printf("%s\n", ability.Description)
		fmt.Printf("⚡ %s\n", ability.Effect)
	}

	// Герой атакує
	attack := b.Hero.Attack()
	if attack.HeadsCut > 0 {
		if attack.Type != NormalAttack {
			fmt.Printf("⚔️ %s використовує %s!\n", b.Hero.Name, attack.Description)
		}
		fmt.Printf("💥 БАМ! %s відрубав %d голів змія!\n", b.Hero.Name, attack.HeadsCut)
	}

	b.Dragon.Heads -= attack.HeadsCut

	if b.Dragon.Heads <= 0 {
		return errors.NewHeroVictoryError()
	}

	if b.Dragon.Heads == b.config.DragonPhobiaNumber {
		return errors.NewDragonPhobiaError()
	}

	// Регенерація з урахуванням вогняного подиху
	regeneratedHeads := b.Dragon.RegenerateHeads(attack.HeadsCut, b.Dragon.Firebreath.Used)

	if regeneratedHeads > 0 {
		b.Dragon.Heads += regeneratedHeads
		if b.Dragon.Firebreath.Used {
			fmt.Printf("🔥 Вогняний подих подвоює регенерацію!\n")
		}
		fmt.Printf("🐉 Ох ні! У змія відросло %d голів!\n", regeneratedHeads)
		fmt.Printf("😱 Тепер у змія %d голів!\n", b.Dragon.Heads)
	} else {
		fmt.Printf("😎 Голови не відросли! У змія залишилось %d голів.\n", b.Dragon.Heads)
	}

	if b.Dragon.Heads >= b.config.MaxDragonHeads {
		return errors.NewDragonVictoryError()
	}

	return nil
}
