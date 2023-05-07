package sqldb

import (
	"encoding/json"

	"github.com/go-sql-driver/mysql"
)

// ////
// Custom SQL Types
// ////

// NullTime aliases msql.NullTime
// It implements a custom MarshallJSON func, which takes the default mysql.NullTime obj of { Time, Valid } into a single
// time or null field in the resulting JSON
type NullTime mysql.NullTime

func (r NullTime) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.Time)
	} else {
		return json.Marshal(nil)
	}
}
