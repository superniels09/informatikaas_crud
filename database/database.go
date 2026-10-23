package database

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"informatikaas_crud/config"
	"informatikaas_crud/sql_queries"
	"strings"

	_ "embed"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Ctx     context.Context
	Queries *sql_queries.Queries

	mysqlDB *sql.DB
}

func New(cfg *config.Config, schema []byte) (*DB, error) {
	db := &DB{
		Ctx: context.Background(),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	mysqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.mysqlDB = mysqlDB

	i := 0
	for part := range bytes.SplitSeq(schema, []byte(";")) {
		partString := string(part)
		if strings.TrimSpace(partString) == "" {
			i++
			continue
		}
		fmt.Println("Executing schema part", i+1)
		if _, err := mysqlDB.ExecContext(db.Ctx, partString); err != nil {
			return nil, err
		}
		i++
	}

	db.Queries = sql_queries.New(mysqlDB)
	return db, nil
}
