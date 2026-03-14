package champ

import (
	"fmt"
	"strings"
)

// PrintCurrentTables печатает текущее состояние пилотов и команд в двух таблицах рядом.
func PrintCurrentTables() {
	fmt.Println("==================================================")
	fmt.Printf(" ТЕКУЩАЯ ТАБЛИЦА ПИЛОТОВ  /  КОМАНД (Осталось ГП: %d, Спринтов: %d)\n", remStages, remSprints)
	fmt.Println("==================================================")

	leftHeader := fmt.Sprintf("| %-12s | %-10s |", "Имя Участника", "Очки")
	rightHeader := fmt.Sprintf("| %-14s | %-10s |", "Команда", "Очки")
	fmt.Printf("%-32s    %-32s\n", leftHeader, rightHeader)

	leftSep := "|--------------|------------|"
	rightSep := "|----------------|------------|"
	fmt.Printf("%-32s    %-32s\n", leftSep, rightSep)

	maxRows := len(initialParticipantsData)
	if len(initialTeamsData) > maxRows {
		maxRows = len(initialTeamsData)
	}

	for i := 0; i < maxRows; i++ {
		left := ""
		right := ""
		if i < len(initialParticipantsData) {
			d := initialParticipantsData[i]
			left = fmt.Sprintf("| %-12s | %-10d |", d.Name, d.Score)
		}
		if i < len(initialTeamsData) {
			t := initialTeamsData[i]
			right = fmt.Sprintf("| %-14s | %-10d |", t.Name, t.Score)
		}
		fmt.Printf("%-32s    %-32s\n", left, right)
	}
	fmt.Println("==================================================")
}

// PrintClinchAnalysis выводит анализ досрочно гарантированных позиций.
func PrintClinchAnalysis() {
	lines := buildDriverClinchLines()
	fmt.Println("==================================================")
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println("==================================================")
}

// PrintConstructorsClinchAnalysis анализирует досрочные позиции в Кубке конструкторов.
func PrintConstructorsClinchAnalysis() {
	lines := buildConstructorsClinchLines()
	fmt.Println("==================================================")
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println("==================================================")
}

func buildDriverClinchLines() []string {
	var lines []string
	lines = append(lines, "       АНАЛИЗ ДОСРОЧНЫХ ПОЗИЦИЙ (CLINCH)")
	lines = append(lines, strings.Repeat("-", 50))

	maxPossibleRemaining := (remStages * maxPoints) + (remSprints * maxSprint)

	foundContest := false
	for i := 0; i < len(initialParticipantsData)-1; i++ {
		leader := initialParticipantsData[i]
		chaser := initialParticipantsData[i+1]
		diff := leader.Score - chaser.Score

		if diff > maxPossibleRemaining {
			lines = append(lines, fmt.Sprintf("ПОЗИЦИЯ %d ГАРАНТИРОВАНА: %-12s (Отрыв: %d)", i+1, leader.Name, diff))
			continue
		}

		needed := maxPossibleRemaining - diff + 1
		lines = append(lines, fmt.Sprintf("БОРЬБА ЗА %d МЕСТО: %s vs %s", i+1, leader.Name, chaser.Name))
		lines = append(lines, fmt.Sprintf("   Осталось разыграть: %d очк.", maxPossibleRemaining))
		lines = append(lines, fmt.Sprintf("   Текущий разрыв: %d очк.", diff))
		lines = append(lines, fmt.Sprintf("   Нужно набрать для гарантии: %d очк.", needed))
		foundContest = true
		break
	}
	if !foundContest {
		lines = append(lines, "Все позиции в чемпионате уже определены!")
	}

	return lines
}

func buildConstructorsClinchLines() []string {
	var lines []string
	lines = append(lines, "   АНАЛИЗ ДОСРОЧНЫХ ПОЗИЦИЙ (КОНСТРУКТОРЫ)")
	lines = append(lines, strings.Repeat("-", 50))

	// максимально команда может набрать очков за этап (две машины)
	maxPossibleRemaining := (remStages * 2 * maxPoints) + (remSprints * 2 * maxSprint)

	foundContest := false
	for i := 0; i < len(initialTeamsData)-1; i++ {
		leader := initialTeamsData[i]
		chaser := initialTeamsData[i+1]
		diff := leader.Score - chaser.Score

		if diff > maxPossibleRemaining {
			lines = append(lines, fmt.Sprintf("ПОЗИЦИЯ %d ГАРАНТИРОВАНА: %-14s (Отрыв: %d)", i+1, leader.Name, diff))
			continue
		}

		needed := maxPossibleRemaining - diff + 1
		lines = append(lines, fmt.Sprintf("БОРЬБА ЗА %d МЕСТО (КОМАНДЫ): %s vs %s", i+1, leader.Name, chaser.Name))
		lines = append(lines, fmt.Sprintf("   Осталось разыграть: %d очк.", maxPossibleRemaining))
		lines = append(lines, fmt.Sprintf("   Текущий разрыв: %d очк.", diff))
		lines = append(lines, fmt.Sprintf("   Нужно набрать для гарантии: %d очк.", needed))
		foundContest = true
		break
	}
	if !foundContest {
		lines = append(lines, "Все позиции в Кубке конструкторов уже определены!")
	}

	return lines
}

// PrintCombinedClinch выводит клинч-анализ пилотов и команд рядом.
func PrintCombinedClinch() {
	driverLines := buildDriverClinchLines()
	teamLines := buildConstructorsClinchLines()

	fmt.Println("==================================================")
	fmt.Println("      CLINCH-ПОЗИЦИИ: ПИЛОТЫ / КОМАНДЫ")
	fmt.Println("==================================================")

	maxLines := len(driverLines)
	if len(teamLines) > maxLines {
		maxLines = len(teamLines)
	}

	for i := 0; i < maxLines; i++ {
		left := ""
		right := ""
		if i < len(driverLines) {
			left = driverLines[i]
		}
		if i < len(teamLines) {
			right = teamLines[i]
		}
		fmt.Printf("%-60s    %-60s\n", left, right)
	}
	fmt.Println("==================================================")
}

