package facades

import (
	"github.com/goravel/inertia/contracts"
	inertiafacades "github.com/goravel/inertia/facades"
)

func Inertia() contracts.Inertia {
	return inertiafacades.Inertia()
}
