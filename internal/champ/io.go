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

// bestPossiblePositions возвращает лучшую возможную итоговую позицию для каждого участника.
// Пилот набирает максимум оставшихся очков; выше него могут быть только те,
// у кого уже сейчас больше очков, чем этот максимум (догнать их невозможно).
func bestPossiblePositions(scores []int, maxRem int) []int {
	best := make([]int, len(scores))
	for i := range scores {
		pMax := scores[i] + maxRem
		guaranteedAhead := 0
		for j := range scores {
			if j == i {
				continue
			}
			if scores[j] > pMax {
				guaranteedAhead++
			}
		}
		best[i] = guaranteedAhead + 1
	}
	return best
}

func partitionClinchLines(lines []string) (header, ceilings, footer []string) {
	inCeilings := false
	for _, line := range lines {
		if strings.Contains(line, "не может подняться") {
			inCeilings = true
			ceilings = append(ceilings, line)
			continue
		}
		if inCeilings {
			footer = append(footer, line)
		} else {
			header = append(header, line)
		}
	}
	return header, ceilings, footer
}

func printSideBySide(left, right []string, width int) {
	maxLines := len(left)
	if len(right) > maxLines {
		maxLines = len(right)
	}
	for i := 0; i < maxLines; i++ {
		l, r := "", ""
		if i < len(left) {
			l = left[i]
		}
		if i < len(right) {
			r = right[i]
		}
		fmt.Printf("%-*s    %-*s\n", width, l, width, r)
	}
}

func buildDriverClinchLines() []string {
	var lines []string
	lines = append(lines, "       АНАЛИЗ ДОСРОЧНЫХ ПОЗИЦИЙ (CLINCH)")
	lines = append(lines, strings.Repeat("-", 50))

	maxPossibleRemaining := (remStages * maxPoints) + (remSprints * maxSprint)

	// Список пилотов, которые уже гарантировали хотя бы одно место.
	// Формат: "NAME(место)".
	var clinchedAny []string

	foundContest := false
	for i := 0; i < len(initialParticipantsData)-1; i++ {
		leader := initialParticipantsData[i]
		chaser := initialParticipantsData[i+1]
		diff := leader.Score - chaser.Score

		if diff > maxPossibleRemaining {
			lines = append(lines, fmt.Sprintf("ПОЗИЦИЯ %d ГАРАНТИРОВАНА: %-12s (Отрыв: %d)", i+1, leader.Name, diff))
			clinchedAny = append(clinchedAny, fmt.Sprintf("%s(%d)", leader.Name, i+1))
			continue
		}

		needed := maxPossibleRemaining - diff + 1
		lines = append(lines, fmt.Sprintf("БОРЬБА ЗА %d МЕСТО: %s vs %s", i+1, leader.Name, chaser.Name))
		lines = append(lines, fmt.Sprintf("   Осталось разыграть: %d очк.", maxPossibleRemaining))
		lines = append(lines, fmt.Sprintf("   Текущий разрыв: %d очк.", diff))
		lines = append(lines, fmt.Sprintf("   Нужно набрать для гарантии: %d очк.", needed))

		scores := make([]int, len(initialParticipantsData))
		for j, p := range initialParticipantsData {
			scores[j] = p.Score
		}
		bestPos := bestPossiblePositions(scores, maxPossibleRemaining)
		for j := i + 2; j < len(initialParticipantsData); j++ {
			p := initialParticipantsData[j]
			pos := bestPos[j]
			if pos > i+1 {
				lines = append(lines, fmt.Sprintf("   %-12s: не может подняться выше %d-го места", p.Name, pos))
			}
		}

		foundContest = true
		break
	}
	if !foundContest {
		lines = append(lines, "Все позиции в чемпионате уже определены!")

		// гарантированным фактически является и последнее место
		// (оно следует из определения всех позиций сверху).
		clinchedAny = clinchedAny[:0]
		for idx, p := range initialParticipantsData {
			clinchedAny = append(clinchedAny, fmt.Sprintf("%s(%d)", p.Name, idx+1))
		}
	} else if len(clinchedAny) == 0 {
		lines = append(lines, "Никто ещё не гарантировал себе место.")
	} else {
		lines = append(lines, "Гарантировали хотя бы одно место: "+strings.Join(clinchedAny, ", "))
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

		scores := make([]int, len(initialTeamsData))
		for j, t := range initialTeamsData {
			scores[j] = t.Score
		}
		bestPos := bestPossiblePositions(scores, maxPossibleRemaining)
		for j := i + 2; j < len(initialTeamsData); j++ {
			t := initialTeamsData[j]
			pos := bestPos[j]
			if pos > i+1 {
				lines = append(lines, fmt.Sprintf("   %-14s: не может подняться выше %d-го места", t.Name, pos))
			}
		}
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

	dHead, dCeil, dFoot := partitionClinchLines(driverLines)
	tHead, tCeil, tFoot := partitionClinchLines(teamLines)

	const colWidth = 60

	fmt.Println("==================================================")
	fmt.Println("      CLINCH-ПОЗИЦИИ: ПИЛОТЫ / КОМАНДЫ")
	fmt.Println("==================================================")

	printSideBySide(dHead, tHead, colWidth)

	if len(dCeil) > 0 || len(tCeil) > 0 {
		fmt.Println()
		paired := len(dCeil)
		if len(tCeil) < paired {
			paired = len(tCeil)
		}
		printSideBySide(
			append([]string{"   Потолок позиций (пилоты):"}, dCeil[:paired]...),
			append([]string{"   Потолок позиций (команды):"}, tCeil[:paired]...),
			colWidth,
		)
		for _, line := range dCeil[paired:] {
			fmt.Println(line)
		}
		for _, line := range tCeil[paired:] {
			fmt.Println(line)
		}
	}

	for _, line := range dFoot {
		fmt.Println(line)
	}
	for _, line := range tFoot {
		fmt.Println(line)
	}
	fmt.Println("==================================================")
}

