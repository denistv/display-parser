package logger

func NewStdWrapper() *STDWrapper{
	return &STDWrapper{}
}

// STDWrapper обертка для логгера из стандартной библиотеки Go
type STDWrapper struct {

}
