package model

type QuestionType string

const (
	RADIO    = "radio"
	CHECKBOX = "checkbox"
	TEXT     = "text"
	NUMBER   = "number"
)

var TypesMap = map[string]QuestionType{
	"radio":    RADIO,
	"checkbox": CHECKBOX,
	"text":     TEXT,
	"number":   NUMBER,
}

type Option struct {
	Option string `json:"option"`
}

type Question struct {
	Title   string       `json:"title"`
	Type    QuestionType `json:"type"`
	Options []Option     `json:"options"`
}

type Quiz struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Author    string     `json:"author"`
	Questions []Question `json:"question"`
	Answers   []string   `json:"answers"`
}
