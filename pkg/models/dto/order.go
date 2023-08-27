package dto

import "fmt"

type OrderDirection string

const (
	OrderByASC  OrderDirection = "ASC"
	OrderByDESC OrderDirection = "DESC"

	OrderDefaultKey = "id"
)

type OrderParam struct {
	Key       string
	Direction OrderDirection
}

func (a OrderParam) Parse() string {
	if a.Key == "" {
		a.Key = OrderDefaultKey
	}

	key := a.Key
	direction := "DESC"
	if a.Direction == OrderByASC {
		direction = "ASC"
	}

	return fmt.Sprintf("%s %s", key, direction)
}
