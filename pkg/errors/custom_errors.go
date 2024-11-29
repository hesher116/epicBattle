package errors

import "fmt"

type HeroVictoryError struct{}

func (e *HeroVictoryError) Error() string {
	return "Богатир переміг! Змій переможений!"
}

func NewHeroVictoryError() error {
	return &HeroVictoryError{}
}

type DragonVictoryError struct{}

func (e *DragonVictoryError) Error() string {
	return "Змій переміг! У нього забагато голів!"
}

func NewDragonVictoryError() error {
	return &DragonVictoryError{}
}

type DragonPhobiaError struct{}

func (e *DragonPhobiaError) Error() string {
	return "У змія фобія на число 7! Він втікає з поля бою!"
}

func NewDragonPhobiaError() error {
	return &DragonPhobiaError{}
}

type UnexpectedError struct {
	Message string
}

func (e *UnexpectedError) Error() string {
	return fmt.Sprintf("Несподівана помилка: %s", e.Message)
}
