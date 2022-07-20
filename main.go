package main

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/laupse/kubegraph/graph"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.String("Kubeconfig", "", "Kubeconfig")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.SendString("Ok")
	})

	app.Get("/api/graph/fields", func(c *fiber.Ctx) error {
		fields := graph.GetFields()
		return c.JSON(fields)
	})

	app.Get("/api/graph/data", func(c *fiber.Ctx) error {
		namespace := c.Query("ns", "default")
		selector := c.Query("selector", "")
		graphData, err := graph.GetData(namespace, selector)
		if err != nil {
			return err
		}
		return c.JSON(graphData)
	})

	app.Listen(":3000")
}
