package internal

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/shoshtari/paroo/internal/configs"
	"go.uber.org/zap"
)

func RunServer(config configs.SectionHTTPServer, logger *zap.Logger) {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	app.Get("/liveness", func(c fiber.Ctx) error {
		logger := logger.With(zap.String("method", "liveness"))
		logger.Debug("new request")

		return c.SendString("{\"alive\": true}")
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(config.Address, fiber.ListenConfig{
		EnablePrefork: true,
	}))
}
