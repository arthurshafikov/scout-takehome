package types

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type (
	TXContextKey string
	TXFunction   func(ctx context.Context) error
)

type GormSlice []string

func (gs *GormSlice) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &gs)
}

func (gs GormSlice) Value() (driver.Value, error) {
	val, err := json.Marshal(gs)
	return string(val), err
}

type WhereConditions struct {
	Where    map[string]interface{}
	WhereAdv map[string]WhereClause
	Between  map[string]BetweenClause
	Or       map[string]interface{}
	Not      map[string]interface{}
	Order    string
}

type WhereClause struct {
	Operator
	Value interface{}
}

type BetweenClause struct {
	Min    any
	Max    any
	Strict bool
}

type Operator string

const (
	OperatorEqual       = "="
	OperatorMore        = ">"
	OperatorMoreOrEqual = ">="
	OperatorLess        = "<"
	OperatorLessOrEqual = "<="
)

func ApplyWhereConditions(whereConditions *WhereConditions) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if whereConditions.Where != nil || whereConditions.WhereAdv != nil {
			db.Where(whereConditions.Apply(db.Session(&gorm.Session{})))
		}

		if whereConditions.Between != nil {
			for field, clause := range whereConditions.Between {
				firstSign := ">="
				secondSign := "<="

				if clause.Strict {
					firstSign = ">"
					secondSign = "<"
				}

				db.Where(
					fmt.Sprintf(
						"%s %s ? AND %s %s ?",
						field, firstSign,
						field, secondSign,
					),
					clause.Min,
					clause.Max,
				)
			}
		}

		if whereConditions.Or != nil {
			db.Or(whereConditions.Or)
		}
		if whereConditions.Not != nil {
			db.Not(whereConditions.Not)
		}
		if whereConditions.Order != "" {
			db.Order(whereConditions.Order)
		}

		return db
	}
}

func ApplyPagination(pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination != nil {
			db.Limit(pagination.GetLimit()).
				Offset(pagination.GetOffset())
		}

		return db
	}
}

func (wc *WhereConditions) Apply(db *gorm.DB) *gorm.DB {
	for key, clause := range wc.WhereAdv {
		db = db.Where(
			fmt.Sprintf("%s %s ?", key, clause.Operator),
			clause.Value,
		)
	}

	return db.Where(wc.Where)
}
