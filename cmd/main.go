package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/stepanov-ds/f1-calc/internal/champ"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := champ.LoadData("data.json"); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	champ.PrintCurrentTables()

	champ.PrintCombinedClinch()

	champ.RunCombinedSimulations()
}
