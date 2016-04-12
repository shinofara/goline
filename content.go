package goline

const (
	TypeText     = 1
	TypeImage    = 2
	TypeVideo    = 3
	TypeAudio    = 4
	TypeLocation = 7
	TypeSticker  = 8
	TypeContact  = 10

	ToType = 1
)

type Content struct {
	ID          string   `json:"id"`
	ContentType int      `json:"contentType"`
	From        string   `json:"from"`
	CreatedTime int      `json:"createdTime"`
	To          []string `json:"to"`
	ToType      int      `json:"toType"`
	Text        string   `json:"text"`
}

func NewToContent(text string) *Content {
	return &Content{
		ContentType: TypeText,
		ToType:      ToType,
		Text:        text,
	}
}
