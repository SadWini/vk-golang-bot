package structs

type Attachment struct {
	AttachType      string
	Attach          string
	Fwd             string
	From            int
	Geo             int
	GeoProvider     int
	Title           string
	AttachProductID int
	AttachPhoto     string
	AttachTitle     string
	AttachDesc      string
	AttachURL       string
	Emoji           bool
	FromAdmin       int
	SourceAct       string
	SourceMid       int
}
