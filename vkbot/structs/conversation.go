package structs

type ConversationInfo struct {
	Peer struct {
		ID      int
		Type    string
		LocalID int `json:"local_id"`
	}
	InRead        int `json:"in_read"`
	OutRead       int `json:"out_read"`
	LastMessageID int `json:"last_message_id"`
	CanWrite      struct {
		Allowed bool
	} `json:"can_write"`
	ChatSettings struct {
		Title        string
		MembersCount int `json:"members_count"`
		State        string
		ActiveIDs    []int `json:"active_ids"`
		ACL          struct {
			CanInvite           bool `json:"can_invite"`
			CanChangeInfo       bool `json:"can_change_info"`
			CanChangePin        bool `json:"can_change_pin"`
			CanPromoteUsers     bool `json:"can_promote_users"`
			CanSeeInviteLink    bool `json:"can_see_invite_link"`
			CanChangeInviteLink bool `json:"can_change_invite_link"`
		}
		IsGroupChannel bool `json:"is_group_channel"`
		OwnerID        int  `json:"owner_id"`
	} `json:"chat_settings"`
}
type ConversationsResponse struct {
	Response struct {
		Items    []ConversationInfo
		Profiles []UserProfile
	}
	Error *VKError
}
