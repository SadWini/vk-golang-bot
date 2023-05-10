package structs

type UserProfile struct {
	ID              int
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	IsClosed        bool   `json:"is_closed"`
	CanAccessClosed bool   `json:"can_access_closed"`
	Sex             int
	ScreenName      string `json:"screen_name"`
	BDate           string `json:"bdate"`
	Photo           string
	Online          int
	City            Geo
	Country         Geo
}

type GroupProfile struct {
	ID       int
	Name     string
	IsClosed int `json:"is_closed"`
	Type     string
	Photo50  string
	Photo100 string
	Photo200 string
}
