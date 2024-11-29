package utils

import (
	"fmt"
	"math/rand"
)

var (
	adjectives1 = []string{"Найкращий", "Найжвавіший", "Найгірший", "Смердючий", "Сміливий", "Найхоробріший"}
	adjectives2 = []string{"сивий", "румяний", "веселий", "випивший", "бородатий", "кремезний"}
	nouns       = []string{"Камінь", "Степан", "Кінь", "Олександр", "Богатир", "Молодець", "Козак"}
)

func GenerateHeroName() string {
	adj1 := adjectives1[rand.Intn(len(adjectives1))]
	adj2 := adjectives2[rand.Intn(len(adjectives2))]
	noun := nouns[rand.Intn(len(nouns))]
	return fmt.Sprintf("%s %s %s", adj1, adj2, noun)
}
