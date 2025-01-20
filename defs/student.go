package defs

type TotalScore map[string]string

type Student struct {
	ID          string
	Name        string
	ClassName   string
	Institution string
	Scores      TotalScore
}
