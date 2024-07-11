package server

import (
	"encoding/json"
	"gwi-test/domain"
	"gwi-test/ports"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	AssetHandler ports.AssetHandler
}

func NewServer(assetsHandler ports.AssetHandler) *Server {
	return &Server{
		AssetHandler: assetsHandler,
	}
}

func (server *Server) Initialize() {
	app := fiber.New()
	app.Use(cors.New())

	assets := app.Group("/assets")
	assets.Get("/:userId", server.GetAssetsById)
	assets.Get("/objects/:userId", server.GetAssetsObjects)
	assets.Delete("/:assetId", server.DeleteAsset)
	assets.Post("/", server.AddAsset)
	assets.Put("/", server.UpdateAsset)

	log.Fatal(app.Listen(":3000"))
}

func (server *Server) GetAssetsById(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("userId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid ID")
	}
	assets, err := server.AssetHandler.GetAssetsById(uint(userId))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(assets)
}

func (server *Server) GetAssetsObjects(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("userId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid ID")
	}
	assets, err := server.AssetHandler.GetAssetsObjects(uint(userId))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(assets)
}

func (server *Server) DeleteAsset(c *fiber.Ctx) error {
	assetId, err := c.ParamsInt("assetId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid ID")
	}
	err = server.AssetHandler.DeleteAsset(uint(assetId))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON("Asset Deleted")
}

func (server *Server) AddAsset(c *fiber.Ctx) error {
	asset := &domain.Asset{}
	err := json.Unmarshal(c.Body(), asset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid Body")
	}
	err = server.AssetHandler.AddAsset(asset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON("Asset Added")
}

func (server *Server) UpdateAsset(c *fiber.Ctx) error {
	asset := &domain.Asset{}
	err := json.Unmarshal(c.Body(), asset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid Body")
	}
	err = server.AssetHandler.UpdateAsset(asset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON("Asset Updated")
}
