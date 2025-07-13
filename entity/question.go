package entity

type Question struct {
	ID                uint
	QuestionTitle     string
	PossibleQuestions []string
	CorrectAnswer     string
	Difficulty        DifficultyRange
	CategoryID        uint
}

type DifficultyRange uint8

const (
	DifficultyRangeEasy DifficultyRange = iota + 1
	DifficultyRangeMedium
	DifficultyRangeHard
)

func (r DifficultyRange) IsValid() bool {
	if r >= DifficultyRangeEasy || r <= DifficultyRangeHard {
		return true
	}
	return false
}
