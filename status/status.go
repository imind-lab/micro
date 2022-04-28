package status

type Code int

func (c Code) String() string {
	return Errors[c]
}

const Success Code = 10000

const (
	Throttled Code = iota + 10101
	CircuitBreak
	Degradation
	GRPCCanceled
	GRPCUnknown
	GRPCInvalidArgument
	GRPCDeadlineExceeded
	GRPCNotFound
	GRPCAlreadyExists
	GRPCPermissionDenied
	GRPCResourceExhausted
	GRPCFailedPrecondition
	GRPCAborted
	GRPCOutOfRange
	GRPCUnimplemented
	GRPCInternal
	GRPCUnavailable
	GRPCDataLoss
	GRPCUnauthenticated
)

const (
	MissingParams Code = iota + 10201
	InvalidParams
	InvalidCredentials
	InvalidSignature
	InvalidCharactersInMessage
	InvalidLanguageType
	InvalidUrl
	InvalidIPAddress
	InvalidJSONFormat
	InvalidMessage
)

const (
	RESTNotEnabled Code = iota + 10301
	PullDataFailed
	SendToQueueFailed
	UploadFileFailed
	CommunicationFailed
)

const (
	DBConnectFailed Code = iota + 10401
	DBQueryFailed
	DBCreateFailed
	DBUpdateFailed
	DBDeleteFailed
	DBNoRowsAffected
	DBRecordNotFound
	DBDeadlineExceeded
)

const (
	RedisDeadlineExceeded Code = iota + 10421
	RedisDataNotExist
)

const (
	KafkaDeadlineExceeded Code = iota + 10431
)

const (
	InternalError Code = iota + 10901
	SystemError
	UnknownError
)

var Errors = map[Code]string{
	Success: "Success",

	Throttled:              "Throttled",
	CircuitBreak:           "CircuitBreak",
	Degradation:            "Degradation",
	GRPCCanceled:           "GRPCCanceled",
	GRPCUnknown:            "GRPCUnknown",
	GRPCInvalidArgument:    "GRPCInvalidArgument",
	GRPCDeadlineExceeded:   "GRPCRequestTimeout",
	GRPCNotFound:           "GRPCNotFound",
	GRPCAlreadyExists:      "GRPCAlreadyExists",
	GRPCPermissionDenied:   "GRPCPermissionDenied",
	GRPCResourceExhausted:  "GRPCResourceExhausted",
	GRPCFailedPrecondition: "GRPCFailedPrecondition",
	GRPCAborted:            "GRPCAborted",
	GRPCOutOfRange:         "GRPCOutOfRange",
	GRPCUnimplemented:      "GRPCUnimplemented",
	GRPCInternal:           "GRPCInternal",
	GRPCUnavailable:        "GRPCUnavailable",
	GRPCDataLoss:           "GRPCDataLoss",
	GRPCUnauthenticated:    "GRPCUnauthenticated",

	MissingParams:              "Missing params",
	InvalidParams:              "Invalid params",
	InvalidCredentials:         "Invalid credentials",
	InvalidSignature:           "Invalid signature",
	InvalidCharactersInMessage: "Invalid characters in message",
	InvalidLanguageType:        "Invalid language type",
	InvalidUrl:                 "Invalid URL",
	InvalidIPAddress:           "Invalid IP address",
	InvalidJSONFormat:          "Invalid JSON format",
	InvalidMessage:             "Invalid message",

	RESTNotEnabled:      "Account not enabled for REST",
	PullDataFailed:      "Pull data failed",
	SendToQueueFailed:   "Send to queue failed",
	UploadFileFailed:    "Upload file failed",
	CommunicationFailed: "Communication failed",

	DBConnectFailed:    "DB connect failed",
	DBQueryFailed:      "DB query failed",
	DBCreateFailed:     "DB creat failed",
	DBUpdateFailed:     "DB save failed",
	DBDeleteFailed:     "DB delete failed",
	DBNoRowsAffected:   "DB no rows affected",
	DBRecordNotFound:   "DB record not found",
	DBDeadlineExceeded: "DB deadline exceeded",

	RedisDeadlineExceeded: "Redis deadline exceeded",
	RedisDataNotExist:     "RedisDataNotExist",

	KafkaDeadlineExceeded: "Kafka deadline exceeded",

	InternalError: "Internal error",
	SystemError:   "System error",
	UnknownError:  "Unknown error",
}
