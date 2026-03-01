package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Seller struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"` // 'chain', 'independent', 'online'
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func main() {
	db, err := sql.Open("postgres", "postgres://etin:etin@localhost:5432/ukoni?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sellers_test (id uuid DEFAULT gen_random_uuid() PRIMARY KEY, name text, type text, created_at timestamp DEFAULT now(), deleted_at timestamp)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO sellers_test (name, type) VALUES ('test', 'chain')")
	if err != nil {
		panic(err)
	}

	query := `SELECT id, name, type, created_at, deleted_at FROM sellers_test WHERE deleted_at IS NULL ORDER BY name ASC`
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	sellers := []*Seller{}
	for rows.Next() {
		var s Seller
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Type,
			&s.CreatedAt,
			&s.DeletedAt,
		)
		if err != nil {
			fmt.Println("Scan error:", err)
		} else {
			fmt.Println("Scanned successfully", s)
		}
		sellers = append(sellers, &s)
	}

	db.Exec("DROP TABLE sellers_test")
}
