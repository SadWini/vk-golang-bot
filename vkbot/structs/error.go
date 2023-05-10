package api

type VKError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	//	RequestParams
}
type ErrorResponse struct {
	Error *VKError
}

type ResponseError struct {
	err     error
	content string
}

func (err ResponseError) Error() string {
	return err.err.Error()
}

func (err ResponseError) Content() string {
	return err.content
}

func (err VKError) Error() string {
	return "vk: " + err.ErrorMsg
}
