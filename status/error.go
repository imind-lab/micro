package status

import (
    "errors"
    "fmt"
)

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
    return fmt.Sprintf("code = %d, message = %v", err.Code, err.Msg)
}

// ResponseOK 处理成功
var ResponseOK = New(Success)

// 微服务gRPC相关error
var ErrThrottled = New(Throttled)
var ErrCircuitBreak = New(CircuitBreak)
var ErrDegradation = New(Degradation)
var ErrGRPCCanceled = New(GRPCCanceled)
var ErrGRPCUnknown = New(GRPCUnknown)
var ErrGRPCInvalidArgument = New(GRPCInvalidArgument)
var ErrGRPCDeadlineExceeded = New(GRPCDeadlineExceeded)
var ErrGRPCNotFound = New(GRPCNotFound)
var ErrGRPCAlreadyExists = New(GRPCAlreadyExists)
var ErrGRPCPermissionDenied = New(GRPCPermissionDenied)
var ErrGRPCResourceExhausted = New(GRPCResourceExhausted)
var ErrGRPCFailedPrecondition = New(GRPCFailedPrecondition)
var ErrGRPCAborted = New(GRPCAborted)
var ErrGRPCOutOfRange = New(GRPCOutOfRange)
var ErrGRPCUnimplemented = New(GRPCUnimplemented)
var ErrGRPCInternal = New(GRPCInternal)
var ErrGRPCUnavailable = New(GRPCUnavailable)
var ErrGRPCDataLoss = New(GRPCDataLoss)
var ErrGRPCUnauthenticated = New(GRPCUnauthenticated)

// 中间件相关error
var ErrDBConnect = New(DBConnectFailed)
var ErrDBQuery = New(DBQueryFailed)
var ErrDBCreate = New(DBCreateFailed)
var ErrDBUpdate = New(DBUpdateFailed)
var ErrDBDelete = New(DBDeleteFailed)
var ErrNoRowsAffected = New(DBNoRowsAffected)
var ErrRecordNotFound = New(DBRecordNotFound)
var ErrDBDeadlineExceeded = New(DBDeadlineExceeded)
var ErrRedisDeadlineExceeded = New(RedisDeadlineExceeded)
var ErrRedisDataNotExist = New(RedisDataNotExist)
var ErrKafkaDeadlineExceeded = New(KafkaDeadlineExceeded)

// ErrorInvalidUrl 业务逻辑相关问题
var ErrInvalidUrl = New(InvalidUrl)
var ErrPullDataFailed = New(PullDataFailed)
var ErrInvalidParams = New(InvalidParams)
var ErrUploadFileFailed = New(UploadFileFailed)

// ErrorSystem 严重错误
var ErrSystem = New(SystemError)

func AsError(err error) Error {
    if err == nil {
        return ResponseOK
    }
    var e Error
    if errors.Is(err, &e) {
        return e
    }
    return New(UnknownError)
}
