package migrations

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

func MustMigrate(db *sqlx.DB) {
	source := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			Migration01(),
		},
	}

	n, err := migrate.ExecMax(db.DB, "postgres", source, migrate.Up, 0)

	if err != nil {
		panic(err)
	}

	if n > 0 {
		fmt.Printf("Applied %d migrations!\n", n)
	}
}
