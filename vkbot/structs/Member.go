package structs

type MemberItem struct {
	MemberID  int  `json:"member_id"`
	JoinDate  int  `json:"join_date"`
	IsOwner   bool `json:"is_owner"`
	IsAdmin   bool `json:"is_admin"`
	InvitedBy int  `json:"invited_by"`
}

type VKMembers struct {
	Items    []MemberItem
	Profiles []UserProfile
	Groups   []GroupProfile
}
