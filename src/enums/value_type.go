package enums

import "database/sql/driver"

type ValueType string

const (
	TEXT ValueType = "text"
	HTML ValueType = "html"
	JSON ValueType = "json"
)

func (vt *ValueType) Scan(value interface{}) error {
	*vt = ValueType(value.(string))
	return nil
}

func (vt ValueType) Value() (driver.Value, error) {
	return string(vt), nil
}

func (vt ValueType) String() string {
	return string(vt)
}
