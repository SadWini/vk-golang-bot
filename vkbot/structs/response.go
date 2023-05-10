package structs

type UsersResponse struct {
	Response VKUsers
	Error    *VKError
}

type MembersResponse struct {
	Response VKMembers
	Error    *VKError
}

type SimpleResponse struct {
	Response int
	Error    *VKError
}

type ChatInfo struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Kicked  int    `json:"kicked"`
	AdminID int    `json:"admin_id"`
	Users   VKUsers
}

type ChatInfoResponse struct {
	Response ChatInfo
	Error    *VKError
}

type FailResponse struct {
	Failed     int
	Ts         int
	MinVersion int `json:"min_version"`
	MaxVersion int `json:"max_version"`
}

type GroupFailResponse struct {
	Failed     int
	Ts         string
	MinVersion int `json:"min_version"`
	MaxVersion int `json:"max_version"`
}
