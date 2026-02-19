package v1

import (
	"worklayer/internal/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Routes interface {
	Register()
}

type routes struct {
	app *fiber.App
	cfg *config.Config
	db  *gorm.DB
}

func New(app *fiber.App, config *config.Config, db *gorm.DB) Routes {
	return &routes{
		app: app,
		cfg: config,
		db:  db,
	}
}

func (r *routes) Register() {
	deps := r.buildDependencies()

	api := r.app.Group("/api/v1")

	r.registerHealthRoutes(api, deps)       // Health
	r.registerAuthRoutes(api, deps)         // Auth
	r.registerUserRoutes(api, deps)         // User
	r.registerOrganizationRoutes(api, deps) // Organization
}
