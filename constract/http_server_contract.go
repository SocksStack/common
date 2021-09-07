package constract

type Callback func()

type IHttpConfig interface {}

type Route interface{}
type Routes []Route

type Middleware interface{}
type Middlewares []Middleware

type IServer interface {
	Run()
	Route(route ...Route) IServer
	Use(middleware ...Middleware) IServer
	SetConfig(config IHttpConfig) IServer
}
