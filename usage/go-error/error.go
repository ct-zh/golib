package go_error

// https://gist.github.com/prathabk/744367cbfc70435c56956f650612d64b

// GoError 接口
type Error interface {
	error

	Code() string
	Message() string
	Causes() []error
	Cause() error
	Data() map[string]interface{}
	String() string
	ResponseType() ResponseErrType
	SetResponseType(errType ResponseErrType) Error
	Component() ErrComponent
	SetComponent(component ErrComponent) Error
	Retryable() bool
	SetRetryable(bool) Error
}

type GoError struct {
	error
	Code         string                 // 错误码
	Data         map[string]interface{} // 上下文数据
	Causes       []error                // 错误堆栈
	Component    ErrComponent           // 标记组件，用于识别error发生在哪一层
	ResponseType ResponseErrType        // 响应类型
	Retryable    bool                   // 重试
}

type ErrComponent string

const (
	ErrService ErrComponent = "service"
	ErrRepo    ErrComponent = "repository"
	ErrLib     ErrComponent = "library"
)

type ResponseErrType string

const (
	BadRequest    ResponseErrType = "BadRequest"
	Forbidden     ResponseErrType = "Forbidden"
	NotFound      ResponseErrType = "NotFound"
	AlreadyExists ResponseErrType = "AlreadyExists"
)
