package champ

import (
	"encoding/json"
	"os"
	"sort"
)

// LoadData загружает историю чемпионата и подготавливает стартовые данные.
func LoadData(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var cfg config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return err
	}

	doneStages, doneSprints := 0, 0
	tempMap := make(map[string]int)    // очки пилотов
	teamTemp := make(map[string]int)   // очки команд
	driverTeams = make(map[string]string)

	for _, gp := range cfg.History {
		if gp.IsSprint {
			doneSprints++
		} else {
			doneStages++
		}
		currMap := pointsMap
		if gp.IsSprint {
			currMap = pointsMapSprint
		}
		for pos, p := range gp.Participants {
			// позиции считаем по порядку появления (1-based)
			place := pos + 1
			points := currMap[place]
			tempMap[p.Name] += points
			if p.Team != "" {
				driverTeams[p.Name] = p.Team
				teamTemp[p.Team] += points
			}
		}
	}

	// коэффициенты по последним 5 событиям (ГП + спринты)
	coefficients = make(map[string]float64)
	if len(cfg.History) > 0 {
		startIdx := 0
		if len(cfg.History) > 5 {
			startIdx = len(cfg.History) - 5
		}
		for _, gp := range cfg.History[startIdx:] {
			for pos, p := range gp.Participants {
				// коэффициент = сумма мест за последние 5 событий
				place := pos + 1
				coefficients[p.Name] += float64(place)
			}
		}
	}

	remStages = cfg.TotalStages - doneStages
	remSprints = cfg.TotalSprints - doneSprints

	initialParticipantsData = nil
	for name, score := range tempMap {
		initialParticipantsData = append(initialParticipantsData, initialData{Name: name, Score: score})
	}
	sort.Slice(initialParticipantsData, func(i, j int) bool {
		return initialParticipantsData[i].Score > initialParticipantsData[j].Score
	})

	initialTeamsData = nil
	for name, score := range teamTemp {
		initialTeamsData = append(initialTeamsData, teamData{Name: name, Score: score})
	}
	sort.Slice(initialTeamsData, func(i, j int) bool {
		return initialTeamsData[i].Score > initialTeamsData[j].Score
	})
	return nil
}

