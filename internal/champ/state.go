package champ

// configuration and shared state

const (
	numSimulations = 1000000
	numWorkers     = 16
	maxPoints      = 25
	maxSprint      = 8
)

var (
	initialParticipantsData []initialData
	initialTeamsData        []teamData
	remStages, remSprints   int
	pointsMap               = map[int]int{1: 25, 2: 18, 3: 15, 4: 12, 5: 10, 6: 8, 7: 6, 8: 4, 9: 2, 10: 1}
	pointsMapSprint         = map[int]int{1: 8, 2: 7, 3: 6, 4: 5, 5: 4, 6: 3, 7: 2, 8: 1}
	coefficients            map[string]float64
	driverTeams             map[string]string
)

