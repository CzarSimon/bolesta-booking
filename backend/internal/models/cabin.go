package models

import (
	"fmt"
	"time"
)

// Cabin house that is bookable
type Cabin struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c Cabin) String() string {
	return fmt.Sprintf("Cabin(id=%s, name=%s, createdAt=%v, updatedAt=%v)", c.ID, c.Name, c.CreatedAt, c.UpdatedAt)
}
