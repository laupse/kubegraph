package ports

type HttpHandler interface {
	Run(address string)
	SetupRoutes()
}
