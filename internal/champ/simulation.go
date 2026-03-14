package champ

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

func runSingleSimulation() string {
	participants := make([]participant, len(initialParticipantsData))
	for i, d := range initialParticipantsData {
		participants[i] = participant{
			Name:          d.Name,
			TotalScore:    d.Score,
			StartingScore: d.Score,
			Coefficient:   coefficients[d.Name],
		}
	}

	// Симуляция ГП
	for s := 0; s < remStages; s++ {
		indices := rand.Perm(len(participants))
		for place, idx := range indices {
			participants[idx].TotalScore += pointsMap[place+1]
		}
	}
	// Симуляция Спринтов
	for s := 0; s < remSprints; s++ {
		indices := rand.Perm(len(participants))
		for place, idx := range indices {
			participants[idx].TotalScore += pointsMapSprint[place+1]
		}
	}

	sort.Slice(participants, func(i, j int) bool {
		if participants[i].TotalScore != participants[j].TotalScore {
			return participants[i].TotalScore > participants[j].TotalScore
		}
		return participants[i].Name < participants[j].Name
	})
	return participants[0].Name
}

func runSingleConstructorsSimulation() string {
	participants := make([]participant, len(initialParticipantsData))
	for i, d := range initialParticipantsData {
		participants[i] = participant{
			Name:          d.Name,
			TotalScore:    d.Score,
			StartingScore: d.Score,
			Coefficient:   coefficients[d.Name],
		}
	}

	teamScores := make(map[string]int)
	for _, t := range initialTeamsData {
		teamScores[t.Name] = t.Score
	}

	// Симуляция ГП
	for s := 0; s < remStages; s++ {
		indices := rand.Perm(len(participants))
		for place, idx := range indices {
			points := pointsMap[place+1]
			participants[idx].TotalScore += points
			if team, ok := driverTeams[participants[idx].Name]; ok {
				teamScores[team] += points
			}
		}
	}
	// Симуляция Спринтов
	for s := 0; s < remSprints; s++ {
		indices := rand.Perm(len(participants))
		for place, idx := range indices {
			points := pointsMapSprint[place+1]
			participants[idx].TotalScore += points
			if team, ok := driverTeams[participants[idx].Name]; ok {
				teamScores[team] += points
			}
		}
	}

	type teamResult struct {
		Name  string
		Score int
	}

	var results []teamResult
	for name, score := range teamScores {
		results = append(results, teamResult{Name: name, Score: score})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Name < results[j].Name
	})

	if len(results) == 0 {
		return ""
	}
	return results[0].Name
}

// RunSimulations запускает серию симуляций и печатает финальную таблицу.
func RunSimulations() {
	stats := simulateDrivers()

	fmt.Println("\n==================================================")
	fmt.Println("       СВОДНАЯ ТАБЛИЦА ПРОЦЕНТА ПОБЕД")
	fmt.Println("==================================================")
	fmt.Printf("| %-12s | %-8s | %-10s |\n", "Имя Участника", "Побед", "Процент")
	fmt.Println("|--------------|----------|------------|")

	for _, s := range stats {
		if s.WinPercentage > 0 {
			fmt.Printf("| %-12s | %-8d | %9.4f%% |\n", s.Name, s.Wins, s.WinPercentage)
		}
	}
	fmt.Println("==================================================")
}

// RunConstructorsSimulations запускает серию симуляций и печатает шансы команд.
func RunConstructorsSimulations() {
	stats := simulateConstructors()

	fmt.Println("\n==================================================")
	fmt.Println("    СВОДНАЯ ТАБЛИЦА ПРОЦЕНТА ПОБЕД КОМАНД")
	fmt.Println("==================================================")
	fmt.Printf("| %-14s | %-8s | %-10s |\n", "Команда", "Побед", "Процент")
	fmt.Println("|----------------|----------|------------|")

	for _, s := range stats {
		if s.WinPercentage > 0 {
			fmt.Printf("| %-14s | %-8d | %9.4f%% |\n", s.Name, s.Wins, s.WinPercentage)
		}
	}
	fmt.Println("==================================================")
}

// RunCombinedSimulations выполняет симуляции и печатает результаты пилотов и команд рядом.
func RunCombinedSimulations() {
	fmt.Printf("Запуск %d симуляций для пилотов и команд...\n", numSimulations)

	driverStats := simulateDrivers()
	teamStats := simulateConstructors()

	fmt.Println("\n==================================================")
	fmt.Println("  СВОДНАЯ ТАБЛИЦА ПРОЦЕНТА ПОБЕД: ПИЛОТЫ / КОМАНДЫ")
	fmt.Println("==================================================")

	leftHeader := fmt.Sprintf("| %-12s | %-8s | %-10s |", "Имя Участника", "Побед", "Процент")
	rightHeader := fmt.Sprintf("| %-14s | %-8s | %-10s |", "Команда", "Побед", "Процент")
	fmt.Printf("%-40s    %-40s\n", leftHeader, rightHeader)

	leftSep := "|--------------|----------|------------|"
	rightSep := "|----------------|----------|------------|"
	fmt.Printf("%-40s    %-40s\n", leftSep, rightSep)

	maxRows := len(driverStats)
	if len(teamStats) > maxRows {
		maxRows = len(teamStats)
	}

	for i := 0; i < maxRows; i++ {
		left := ""
		right := ""
		if i < len(driverStats) && driverStats[i].WinPercentage > 0 {
			s := driverStats[i]
			left = fmt.Sprintf("| %-12s | %-8d | %9.4f%% |", s.Name, s.Wins, s.WinPercentage)
		}
		if i < len(teamStats) && teamStats[i].WinPercentage > 0 {
			s := teamStats[i]
			right = fmt.Sprintf("| %-14s | %-8d | %9.4f%% |", s.Name, s.Wins, s.WinPercentage)
		}
		fmt.Printf("%-40s    %-40s\n", left, right)
	}
	fmt.Println("==================================================")
}

func simulateDrivers() []winStat {
	fmt.Printf("Запуск %d симуляций (пилоты) в %d потоках...\n", numSimulations, numWorkers)

	resultsCh := make(chan string, 1000)
	var wg sync.WaitGroup
	var mu sync.Mutex
	winCounts := make(map[string]int)
	counter := 0

	for _, d := range initialParticipantsData {
		winCounts[d.Name] = 0
	}

	go func() {
		for winnerName := range resultsCh {
			mu.Lock()
			winCounts[winnerName]++
			counter++
			mu.Unlock()
		}
	}()

	simsPerWorker := numSimulations / numWorkers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < simsPerWorker; i++ {
				resultsCh <- runSingleSimulation()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Прогресс-бар
	ticker := time.NewTicker(100 * time.Millisecond)
	barWidth := 40
	for range ticker.C {
		mu.Lock()
		curr := counter
		mu.Unlock()

		progress := float64(curr) / float64(numSimulations)
		filled := int(progress * float64(barWidth))
		bar := "[" + strings.Repeat("#", filled) + strings.Repeat(" ", barWidth-filled) + "]"
		fmt.Printf("\r%s %5.2f%% (%d/%d)", bar, progress*100, curr, numSimulations)

		if curr >= numSimulations {
			break
		}
	}
	ticker.Stop()
	fmt.Println("\n\nСимуляция (пилоты) завершена.")

	var stats []winStat
	for name, wins := range winCounts {
		stats = append(stats, winStat{
			Name:          name,
			Wins:          wins,
			WinPercentage: (float64(wins) / numSimulations) * 100.0,
		})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].WinPercentage > stats[j].WinPercentage
	})

	return stats
}

func simulateConstructors() []winStat {
	fmt.Printf("Запуск %d симуляций (команды) в %d потоках...\n", numSimulations, numWorkers)

	resultsCh := make(chan string, 1000)
	var wg sync.WaitGroup
	var mu sync.Mutex
	winCounts := make(map[string]int)
	counter := 0

	for _, t := range initialTeamsData {
		winCounts[t.Name] = 0
	}

	go func() {
		for winnerName := range resultsCh {
			if winnerName == "" {
				continue
			}
			mu.Lock()
			winCounts[winnerName]++
			counter++
			mu.Unlock()
		}
	}()

	simsPerWorker := numSimulations / numWorkers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < simsPerWorker; i++ {
				resultsCh <- runSingleConstructorsSimulation()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Прогресс-бар
	ticker := time.NewTicker(100 * time.Millisecond)
	barWidth := 40
	for range ticker.C {
		mu.Lock()
		curr := counter
		mu.Unlock()

		progress := float64(curr) / float64(numSimulations)
		filled := int(progress * float64(barWidth))
		bar := "[" + strings.Repeat("#", filled) + strings.Repeat(" ", barWidth-filled) + "]"
		fmt.Printf("\r%s %5.2f%% (%d/%d)", bar, progress*100, curr, numSimulations)

		if curr >= numSimulations {
			break
		}
	}
	ticker.Stop()
	fmt.Println("\n\nСимуляция (команды) завершена.")

	var stats []winStat
	for name, wins := range winCounts {
		stats = append(stats, winStat{
			Name:          name,
			Wins:          wins,
			WinPercentage: (float64(wins) / numSimulations) * 100.0,
		})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].WinPercentage > stats[j].WinPercentage
	})

	return stats
}

