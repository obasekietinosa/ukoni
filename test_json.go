package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Seller struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"` // 'chain', 'independent', 'online'
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func main() {
	var sellers []*Seller // Empty slice initialized like in List()

	// What does it output?
	b, _ := json.Marshal(sellers)
	fmt.Println(string(b))

	sellers = make([]*Seller, 0)
	b2, _ := json.Marshal(sellers)
	fmt.Println(string(b2))
}
