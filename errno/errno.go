package errno

var (
	//OK Common errors
	OK                    = &Errno{Code: 0, Message: "OK"}
	InternalServerError   = &Errno{Code: 500, Message: "InternalServerError"}
	Unkonwn               = &Errno{Code: 404, Message: "Unkonwn"}
	RequestBodyParseError = &Errno{Code: 501, Message: "RequestBodyParseError"}
)

//Errno the error number for err
type Errno struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err Errno) Error() string {
	return err.Message
}

//DecodeErr decode the error message and code
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}
	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Message

	}
	return InternalServerError.Code, err.Error()

}
