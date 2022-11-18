package http

import (
	json "github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/laupse/kubegraph/application/entity"
	"github.com/laupse/kubegraph/application/services"
)

type FiberHandler struct {
	app          *fiber.App
	graphService *services.GraphService
}

func NewFiberHandler(graphService *services.GraphService) *FiberHandler {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	return &FiberHandler{
		app:          app,
		graphService: graphService,
	}
}

func (f *FiberHandler) Run(address string) {
	f.app.Listen(address)
}

func (f *FiberHandler) SetupRoutes() {
	f.app.Use(pprof.New())
	f.app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.SendString("Ok")
	})

	f.app.Get("/api/graph/fields", func(c *fiber.Ctx) error {
		return c.JSON(f.graphService.GetFields())
	})

	f.app.Get("/api/graph/data", func(c *fiber.Ctx) error {
		grapData, err := f.DataHandler(c, f.graphService)
		if err != nil {
			return err
		}
		return c.JSON(grapData)
	})
}

func (f *FiberHandler) DataHandler(c *fiber.Ctx, grahpService *services.GraphService) (*entity.GraphData, error) {
	namespace := c.Query("ns", "default")
	selector := c.Query("selector", "")
	graphData, err := f.graphService.GetData(namespace, selector)
	return graphData, err
}
