package structs

type FriendRequests struct {
	Count int
	Items []int
}

type FriendRequestsResponse struct {
	Response FriendRequests
	Error    *VKError
}

type FriendDeleteResponse struct {
	Response map[string]int
	Error    *VKError
}
