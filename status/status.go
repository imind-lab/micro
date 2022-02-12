package status

import (
	"fmt"
)

type Code int

func (c Code) String() string {
	return Errors[c]
}

const (
	Success Code = iota + 1000
	Throttled
	MissingParams
	InvalidParams
	InvalidCredentials
	InternalError
	InvalidMessage
	RESTNotEnabled
	CommunicationFailed
	InvalidSignature
	InvalidCharactersInMessage
	InvalidLanguageType
	InvalidUrl
	InvalidIPAddress
	InvalidHTTPStatusCode
	EmptyReturnValue
	InvalidJSONFormat
	PullDataFailed
	SendToQueueFailed
	DBConnectFailed
	DBQueryFailed
	DBSaveFailed
	DBNoRowsAffected
	RouteNotFound
	RecordNotExist
	SystemError
	OtherErrors
)

var Errors = map[Code]string{
	Success:                    "Success",
	Throttled:                  "Throttled",
	MissingParams:              "Missing params",
	InvalidParams:              "Invalid params",
	InvalidCredentials:         "Invalid credentials",
	InternalError:              "Internal error",
	InvalidMessage:             "Invalid message",
	RESTNotEnabled:             "Account not enabled for REST",
	CommunicationFailed:        "Communication failed",
	InvalidSignature:           "Invalid signature",
	InvalidCharactersInMessage: "Invalid characters in message",
	InvalidLanguageType:        "Invalid language type",
	InvalidUrl:                 "Invalid URL",
	InvalidIPAddress:           "Invalid IP address",
	InvalidHTTPStatusCode:      "Invalid HTTP status code",
	EmptyReturnValue:           "Empty return value",
	InvalidJSONFormat:          "Invalid JSON format",
	PullDataFailed:             "Pull data failed",
	SendToQueueFailed:          "Send to queue failed",
	DBConnectFailed:            "DB connect failed",
	DBQueryFailed:              "DB query failed",
	DBSaveFailed:               "DB Save failed",
	DBNoRowsAffected:           "DB No rows affected",
	RouteNotFound:              "Route not found",
	RecordNotExist:             "Record not found",
	SystemError:                "System error",
	OtherErrors:                "Unknown error",
}

type Error struct {
	Code Code
	Msg  string
}

func New(code Code) Error {
	return Error{
		Code: code,
		Msg:  code.String(),
	}
}

// Error 实现Error接口
func (err Error) Error() string {
	return fmt.Sprintf("code = %v, message = %v", err.Code, err.Msg)
}

// ResponseOK 处理成功
var ResponseOK = New(Success)

// ErrorDBConnect 数据库相关问题
var ErrorDBConnect = New(DBConnectFailed)
var ErrorDBQuery = New(DBQueryFailed)
var ErrorDBSave = New(DBSaveFailed)
var ErrorDBNoRowsAffected = New(DBNoRowsAffected)

// ErrorInvalidUrl 业务逻辑相关问题
var ErrorInvalidUrl = New(InvalidUrl)
var ErrorEmptyReturnValue = New(EmptyReturnValue)
var ErrorPullDataFailed = New(PullDataFailed)
var ErrorRouteNotFound = New(RouteNotFound)
var ErrorInvalidParams = New(InvalidParams)
var ErrorRecordNotExist = New(RecordNotExist)

// ErrorSystem 严重错误
var ErrorSystem = New(SystemError)
