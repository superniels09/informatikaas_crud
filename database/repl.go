package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func (db *DB) StartREPL() error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			fmt.Print("> ")
			continue
		}
		if query == "exit" || query == "quit" {
			break
		}
		runQuery(db.mysqlDB, query)
		fmt.Print("> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "input error:", err)
	}
	return nil
}

func runQuery(mysqlDB *sql.DB, query string) {
	rows, err := mysqlDB.Query(query)
	if err != nil {
		// Not a query (e.g. INSERT/UPDATE/DDL) — try Exec instead
		res, execErr := mysqlDB.Exec(query)
		if execErr != nil {
			fmt.Println("error:", err)
			return
		}
		affected, _ := res.RowsAffected()
		fmt.Printf("OK, %d rows affected\n", affected)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Print header
	fmt.Println(strings.Join(cols, "\t"))

	// Buffers to scan into — sql.RawBytes avoids type-specific structs
	vals := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(cols))
	for i := range vals {
		scanArgs[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			fmt.Println("scan error:", err)
			return
		}
		strs := make([]string, len(vals))
		for i, v := range vals {
			if v == nil {
				strs[i] = "NULL"
			} else {
				strs[i] = string(v)
			}
		}
		fmt.Println(strings.Join(strs, "\t"))
	}
	if err := rows.Err(); err != nil {
		fmt.Println("row iteration error:", err)
	}
}
