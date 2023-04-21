package interfaces

type IServiceProvider interface {
	SelectServer() string
	Close()
}
