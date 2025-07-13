package entity

type Game struct {
	ID          uint
	Category    string
	PlayerIDs   []uint
	QuestionIDs []uint
}

type Player struct {
	ID           uint
	UserID       uint
	Answers      []PlayerAnswer
	GameID       uint
	PointsEarned uint16
}

type PlayerAnswer struct {
	ID       uint
	PlayerID uint
	GameID   uint
	Answer   []AnswerOption
}

type AnswerOption uint8

const (
	AnswerOptionA AnswerOption = iota + 1
	AnswerOptionB
	AnswerOptionC
	AnswerOptionD
)
