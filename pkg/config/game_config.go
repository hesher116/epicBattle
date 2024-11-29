package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

const (
	// Ліміти для голів
	MinValidHeads = 1
	MaxValidHeads = 1000

	// Ліміти для атаки
	MinValidAttack = 1
	MaxValidAttack = 50

	// Ліміти для шансів
	MinChance = 0
	MaxChance = 100

	// Стандартні значення
	DefaultMinHeads     = 80
	DefaultMaxHeads     = 120
	DefaultVictoryHeads = 200
	DefaultMinAttack    = 3
	DefaultMaxAttack    = 7
	DefaultPhobiaNumber = 7
)

type GameConfig struct {
	InitialDragonHeadsMin int // Мінімальна початкова кількість голів змія
	InitialDragonHeadsMax int // Максимальна початкова кількість голів змія
	MaxDragonHeads        int // Кількість голів для перемоги змія
	MinHeroAttack        int // Мінімальна сила звичайної атаки героя
	MaxHeroAttack        int // Максимальна сила звичайної атаки героя
	DragonPhobiaNumber   int // Число-фобія змія
	SpecialAttackChance  int // Шанс спеціальної атаки (%)
	FirebreathChance     int // Шанс вогняного подиху (%)
	StunChance           int // Шанс оглушення (%)
	NoRegenChance        int // Шанс що голова не відросте
	OneHeadChance        int // Шанс на 1 голову
	TwoHeadChance        int // Шанс на 2 голови
}

// ValidateRange перевіряє чи знаходиться значення в допустимому діапазоні
func ValidateRange(name string, value, min, max int) error {
	if value < min || value > max {
		return fmt.Errorf("%s має бути між %d та %d, отримано: %d", name, min, max, value)
	}
	return nil
}

// ValidateConfig перевіряє всю конфігурацію на валідність
func (c *GameConfig) Validate() error {
	// Валідація голів дракона
	if err := ValidateRange("мінімальна кількість голів", c.InitialDragonHeadsMin, MinValidHeads, MaxValidHeads); err != nil {
		return err
	}
	if err := ValidateRange("максимальна кількість голів", c.InitialDragonHeadsMax, MinValidHeads, MaxValidHeads); err != nil {
		return err
	}
	if c.InitialDragonHeadsMax <= c.InitialDragonHeadsMin {
		return fmt.Errorf("максимальна кількість голів (%d) має бути більше мінімальної (%d)", 
			c.InitialDragonHeadsMax, c.InitialDragonHeadsMin)
	}

	// Валідація атаки героя
	if err := ValidateRange("мінімальна атака героя", c.MinHeroAttack, MinValidAttack, MaxValidAttack); err != nil {
		return err
	}
	if err := ValidateRange("максимальна атака героя", c.MaxHeroAttack, MinValidAttack, MaxValidAttack); err != nil {
		return err
	}
	if c.MaxHeroAttack <= c.MinHeroAttack {
		return fmt.Errorf("максимальна атака (%d) має бути більше мінімальної (%d)", 
			c.MaxHeroAttack, c.MinHeroAttack)
	}

	// Валідація шансів
	chances := map[string]int{
		"спеціальної атаки": c.SpecialAttackChance,
		"вогняного подиху":  c.FirebreathChance,
		"оглушення":         c.StunChance,
		"без регенерації":   c.NoRegenChance,
		"однієї голови":     c.OneHeadChance,
		"двох голів":        c.TwoHeadChance,
	}

	for name, chance := range chances {
		if err := ValidateRange(fmt.Sprintf("шанс %s", name), chance, MinChance, MaxChance); err != nil {
			return err
		}
	}

	// Валідація послідовності шансів регенерації
	if c.OneHeadChance <= c.NoRegenChance {
		return fmt.Errorf("шанс однієї голови (%d) має бути більше шансу без регенерації (%d)", 
			c.OneHeadChance, c.NoRegenChance)
	}
	if c.TwoHeadChance <= c.OneHeadChance {
		return fmt.Errorf("шанс двох голів (%d) має бути більше шансу однієї голови (%d)", 
			c.TwoHeadChance, c.OneHeadChance)
	}

	return nil
}

// LoadConfig завантажує та валідує конфігурацію
func LoadConfig(envFile string) (GameConfig, error) {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Printf("Попередження: Не вдалося завантажити .env файл: %v\n", err)
		}
	}

	config := GameConfig{
		InitialDragonHeadsMin: getEnvInt("DRAGON_MIN_HEADS", DefaultMinHeads),
		InitialDragonHeadsMax: getEnvInt("DRAGON_MAX_HEADS", DefaultMaxHeads),
		MaxDragonHeads:        getEnvInt("DRAGON_VICTORY_HEADS", DefaultVictoryHeads),
		MinHeroAttack:         getEnvInt("HERO_MIN_ATTACK", DefaultMinAttack),
		MaxHeroAttack:         getEnvInt("HERO_MAX_ATTACK", DefaultMaxAttack),
		DragonPhobiaNumber:    getEnvInt("DRAGON_PHOBIA_NUMBER", DefaultPhobiaNumber),
		SpecialAttackChance:   getEnvInt("HERO_SPECIAL_CHANCE", 20),
		FirebreathChance:      getEnvInt("DRAGON_FIREBREATH_CHANCE", 70),
		StunChance:            getEnvInt("DRAGON_STUN_CHANCE", 15),
		NoRegenChance:         getEnvInt("DRAGON_NO_REGEN_CHANCE", 55),
		OneHeadChance:         getEnvInt("DRAGON_ONE_HEAD_CHANCE", 90),
		TwoHeadChance:         getEnvInt("DRAGON_TWO_HEAD_CHANCE", 98),
	}

	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("невалідна конфігурація: %w", err)
	}

	return config, nil
}

func getEnvInt(key string, defaultVal int) int {
	if val, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
