package enums

import "database/sql/driver"

type TerritoryType string

const (
	COUNTRY TerritoryType = "country"
	NUMERIC TerritoryType = "numeric"
)

func (tt *TerritoryType) Scan(value interface{}) error {
	*tt = TerritoryType(value.(string))
	return nil
}

func (tt TerritoryType) Value() (driver.Value, error) {
	return string(tt), nil
}

func (tt TerritoryType) String() string {
	return string(tt)
}
