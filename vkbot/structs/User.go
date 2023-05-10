package structs

import "strings"

type Geo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type User struct {
	ID              int    `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ScreenName      string `json:"screen_name"`
	Photo           string `json:"photo"`
	InvitedBy       int    `json:"invited_by"`
	City            Geo    `json:"city"`
	Country         Geo    `json:"country"`
	Sex             int    `json:"sex"`
	BDate           string `json:"bdate"`
	Photo50         string `json:"photo_50"`
	Photo100        string `json:"photo_100"`
	Status          string `json:"status"`
	About           string `json:"about"`
	Relation        int    `json:"relation"`
	Hidden          int    `json:"hidden"`
	Closed          int    `json:"is_closed"`
	CanAccessClosed int    `json:"can_access_closed"`
	Deactivated     string `json:"deactivated"`
	IsAdmin         bool   `json:"is_admin"`
	IsOwner         bool   `json:"is_owner"`
}
type VKUsers []*User

func (u *User) FullName() string {
	if u != nil {
		return strings.Trim(u.FirstName+" "+u.LastName, " ")
	}
	return ""
}
func (a VKUsers) Len() int           { return len(a) }
func (a VKUsers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VKUsers) Less(i, j int) bool { return a[i].FullName() < a[j].FullName() }
