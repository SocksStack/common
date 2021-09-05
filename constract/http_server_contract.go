package constract

type Callback func()

type IHttpConfig interface {}

type Route interface{}
type Routes []Route

type Middleware interface{}
type Middlewares []Middleware

type IServer interface {
	Run(config IHttpConfig)
	Route(route ...Route)
	Use(middleware ...Middleware)
}
