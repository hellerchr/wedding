package pg

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func NewStatementBuilder(db *sqlx.DB) squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db.DB)
}

func Insert(db *sqlx.DB, table string, arg interface{}) (sql.Result, error) {
	columns := getStructFields("db", arg)
	sColumns := strings.Join(columns, ",")
	sColumnNames := ":" + strings.Join(columns, ",:")

	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", table, sColumns, sColumnNames) //nolint: gosec

	return sqlx.NamedExec(db, sql, arg)
}

func getStructFields(tag string, values interface{}) []string {
	v := reflect.ValueOf(values)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fields := []string{}

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i).Tag.Get(tag)
			if field != "" {
				fields = append(fields, field)
			}
		}

		return fields
	}

	if v.Kind() == reflect.Map {
		for _, keyv := range v.MapKeys() {
			fields = append(fields, keyv.String())
		}

		return fields
	}

	panic(fmt.Errorf("DBFields requires a struct or a map, found: %s", v.Kind().String()))
}
