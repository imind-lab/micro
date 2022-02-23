package status

import (
	"fmt"
	"strings"
)

type Code int

func (c Code) String() string {
	return Errors[c]
}

const (
	Success Code = iota + 1000
	Throttled
	CircuitBroke
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
	CircuitBroke:               "CircuitBroke",
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
	Code    Code
	Message string
}

func NewError(code Code, message ...string) Error {
	msg := code.String()
	if len(message) > 0 {
		msg = strings.Join(message, ", ")
	}
	return Error{
		Code:    code,
		Message: msg,
	}
}

// Error 实现Error接口
func (err Error) Error() string {
	return fmt.Sprintf("code = %v, message = %v", err.Code, err.Message)
}

// ResponseOK 处理成功
var ResponseOK = NewError(Success)

// ErrorDBConnect 数据库相关问题
var ErrorDBConnect = NewError(DBConnectFailed)
var ErrorDBQuery = NewError(DBQueryFailed)
var ErrorDBSave = NewError(DBSaveFailed)
var ErrorDBNoRowsAffected = NewError(DBNoRowsAffected)

// ErrorInvalidUrl 业务逻辑相关问题
var ErrorInvalidUrl = NewError(InvalidUrl)
var ErrorEmptyReturnValue = NewError(EmptyReturnValue)
var ErrorPullDataFailed = NewError(PullDataFailed)
var ErrorRouteNotFound = NewError(RouteNotFound)
var ErrorInvalidParams = NewError(InvalidParams)
var ErrorRecordNotExist = NewError(RecordNotExist)

// ErrorSystem 严重错误
var ErrorSystem = NewError(SystemError)
