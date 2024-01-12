package internal

import "time"

type Product struct {
	Id           int
	Name         string
	Quantity     int
	Code_value   string
	Is_published bool
	Expiration   time.Time
	Price        float64
}
