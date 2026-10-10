package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20261007130631CreatePasswordResetTokensTable struct{}

// Signature The unique signature for the migration.
func (r *M20261007130631CreatePasswordResetTokensTable) Signature() string {
	return "20261007130631_create_password_reset_tokens_table"
}

// Up Run the migrations.
func (r *M20261007130631CreatePasswordResetTokensTable) Up() error {
	if !facades.Schema().HasTable("password_reset_tokens") {
		return facades.Schema().Create("password_reset_tokens", func(table schema.Blueprint) {
			table.String("email")
			table.Unique("email")
			table.String("token_hash", 64)
			table.TimestampTz("expires_at")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261007130631CreatePasswordResetTokensTable) Down() error {
	return facades.Schema().DropIfExists("password_reset_tokens")
}
