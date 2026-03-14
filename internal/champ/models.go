package champ

// domain models

type config struct {
	TotalStages  int         `json:"total_stages"`
	TotalSprints int         `json:"total_sprints"`
	History      []grandPrix `json:"history"`
}

type jsonParticipant struct {
	Name string `json:"name"`
	Team string `json:"team"`
}

type grandPrix struct {
	Name         string            `json:"name"`
	IsSprint     bool              `json:"is_sprint"`
	Participants []jsonParticipant `json:"participants"`
}

type participant struct {
	Name          string
	TotalScore    int
	StartingScore int
	Coefficient   float64
}

type initialData struct {
	Name  string
	Score int
}

type teamData struct {
	Name  string
	Score int
}

type winStat struct {
	Name          string
	Wins          int
	WinPercentage float64
}

