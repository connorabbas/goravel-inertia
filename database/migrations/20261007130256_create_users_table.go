package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20261007130256CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20261007130256CreateUsersTable) Signature() string {
	return "20261007130256_create_users_table"
}

// Up Run the migrations.
func (r *M20261007130256CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		return facades.Schema().Create("users", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("email")
			table.Unique("email")
			table.String("password")
			table.TimestampTz("email_verified_at").Nullable()
			table.String("remember_token", 64).Nullable()
			table.TimestampTz("remember_expires_at").Nullable()
			table.Integer("auth_version").Default(0)
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261007130256CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
