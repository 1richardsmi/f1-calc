package champ

import "testing"

func TestBestPossiblePositions(t *testing.T) {
	t.Run("leader can still win when rivals may score poorly", func(t *testing.T) {
		// HAM 125, ANT 171, RUS 131 — 14 ГП + 3 спринта осталось
		scores := []int{171, 131, 125}
		maxRem := 14*maxPoints + 3*maxSprint // 374

		best := bestPossiblePositions(scores, maxRem)

		if best[0] != 1 {
			t.Fatalf("ANT best=%d, want 1", best[0])
		}
		if best[1] != 1 {
			t.Fatalf("RUS best=%d, want 1", best[1])
		}
		if best[2] != 1 {
			t.Fatalf("HAM best=%d, want 1 (может выиграть, если лидеры не доберут очки)", best[2])
		}
	})

	t.Run("driver eliminated when rival already exceeds their max", func(t *testing.T) {
		// ANT 171, PIA 80 — осталось 2 ГП + 3 спринта (74 очка)
		scores := []int{171, 80}
		maxRem := 2*maxPoints + 3*maxSprint // 74

		best := bestPossiblePositions(scores, maxRem)

		if best[0] != 1 {
			t.Fatalf("ANT best=%d, want 1", best[0])
		}
		if best[1] != 2 {
			t.Fatalf("PIA best=%d, want 2 (171 > 80+74, догнать ANT невозможно)", best[1])
		}
	})
}
