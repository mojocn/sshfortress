package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonArrayString []string

func (o JsonArrayString) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *JsonArrayString) Scan(input interface{}) (err error) {

	switch v := input.(type) {
	case []byte:
		return json.Unmarshal(v, o)
	case string:
		return json.Unmarshal([]byte(v), o)
	default:
		err = fmt.Errorf("unexpected type %T in JsonArraySshFilter", v)
	}
	return err
}
