package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/foundation"

	"goravel/app/http/responses"
	"goravel/config"
	"goravel/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithMigrations(Migrations).
		WithMiddleware(func(middleware configuration.Middleware) {
			middleware.Recover(responses.Recover)
		}).
		WithRouting(func() {
			routes.Web()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
