package category

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (cat *Category) Scan(value interface{}) error {
	catId, ok := value.(int64)
	if !ok {
		return fmt.Errorf("can't scan category of type %v with id %v", reflect.TypeOf(value), value)
	}

	for _, category := range Categories {
		if category.ID == uint(catId) {
			*cat = category
			break
		}
	}

	return nil
}

func (cat Category) Value() (driver.Value, error) {
	for _, category := range Categories {
		if cat.ID == category.ID {
			return int64(cat.ID), nil
		}
	}

	return nil, fmt.Errorf("category with id %d not found", cat.ID)
}

var Categories = []Category{
	{
		ID:   1,
		Name: "CATEGORY_GAMES",
	},
	{
		ID:   2,
		Name: "CATEGORY_BLOG",
	},
	{
		ID:   3,
		Name: "CATEGORY_MUSIC",
	},
	{
		ID:   4,
		Name: "CATEGORY_HUMOR",
	},
	{
		ID:   5,
		Name: "CATEGORY_EDUCATION",
	},
}
