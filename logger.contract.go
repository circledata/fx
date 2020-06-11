package fx

type Logger interface {
	Trace(interface{})
	Debug(interface{})
	Info(interface{})
	Warn(interface{})
	Error(interface{})
}
