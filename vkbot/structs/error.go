package structs

type VKError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	//	RequestParams
}
type ErrorResponse struct {
	Error *VKError
}

type ResponseError struct {
	Err     error
	Content string
}

func (err ResponseError) Error() string {
	return err.Err.Error()
}

func (err ResponseError) getContent() string {
	return err.Content
}

func (err VKError) Error() string {
	return "vk: " + err.ErrorMsg
}
