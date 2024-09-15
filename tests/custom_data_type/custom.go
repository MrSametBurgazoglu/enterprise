package custom_data_type

import (
	"database/sql/driver"
	"fmt"
)

type Custom struct {
	Hello string
	World string
}

func (c Custom) Scan(src any) (err error) {
	switch v := src.(type) {
	case nil:
	case string:
		_, err = fmt.Sscanf(v, "(%q,%q)", &c.Hello, &c.World)
	case []byte:
		_, err = fmt.Sscanf(string(v), "(%q,%q)", &c.Hello, &c.World)
	}
	return
}

func (c Custom) Value() (driver.Value, error) {
	return fmt.Sprintf("(%q,%q)", c.Hello, c.World), nil
}
