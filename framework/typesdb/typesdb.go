package typesdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/DmitryVesenniy/go-rest-framework/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StringNullable
type StringNullable string

func (sn StringNullable) Value() (driver.Value, error) {

	if sn == "" {
		return nil, nil
	}

	return string(sn), nil
}
func (sn StringNullable) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if sn == "" {
		return clause.Expr{
			SQL:  "?",
			Vars: []interface{}{nil},
		}
	}
	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{string(sn)},
	}
}
func (sn *StringNullable) Scan(v interface{}) error {
	if v == nil {
		*sn = StringNullable("")
		return nil
	}
	value := common.ConvertToString(v)
	*sn = StringNullable(value)
	return nil
}

// Uint64Nullable
type UintNullable uint

func (id UintNullable) GormDataType() string {
	return "bigint"
}

func (id UintNullable) Value() (driver.Value, error) {
	if id == 0 {
		return nil, nil
	}

	return uint(id), nil
}

func (id UintNullable) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if id == 0 {
		return clause.Expr{
			SQL:  "?",
			Vars: []interface{}{nil},
		}
	}
	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{uint(id)},
	}
}

func (id *UintNullable) Scan(v interface{}) error {
	value, err := common.ConvertToInt(v)
	if err != nil {
		*id = UintNullable(0)
		return nil
	}
	*id = UintNullable(value)
	return nil
}

// Uint64Nullable
type Uint64Nullable uint64

func (id Uint64Nullable) GormDataType() string {
	return "bigint"
}
func (id Uint64Nullable) Value() (driver.Value, error) {
	if id == 0 {
		return nil, nil
	}

	return uint(id), nil
}
func (id Uint64Nullable) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if id == 0 {
		return clause.Expr{
			SQL:  "?",
			Vars: []interface{}{nil},
		}
	}
	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{uint64(id)},
	}
}
func (id *Uint64Nullable) Scan(v interface{}) error {
	value, err := common.ConvertToInt(v)
	if err != nil {
		*id = Uint64Nullable(0)
		return nil
	}
	*id = Uint64Nullable(value)
	return nil
}

// BitBool
type BitBool bool

func (b BitBool) GormDataType() string {
	return "bit"
}

func (b BitBool) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	var result byte
	if b {
		result = 1
	} else {
		result = 0
	}
	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{(result)},
	}
}

func (b BitBool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	} else {
		return []byte{0}, nil
	}
}

func (b *BitBool) Scan(src interface{}) error {
	v, ok := src.([]byte)
	if !ok {
		return errors.New("bad []byte type assertion")
	}
	*b = v[0] == 1
	return nil
}

// UUID1
type UUID1 string

func (uuid UUID1) GormDataType() string {
	return "uuid"
}
func (uuid UUID1) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if uuid == "" {
		return clause.Expr{
			SQL:  "?",
			Vars: []interface{}{nil},
		}
	}
	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{string(uuid)},
	}
}
func (uuid *UUID1) Scan(v interface{}) error {
	res, ok := v.([]byte)
	if ok {
		*uuid = UUID1(res)
	} else {
		*uuid = ""
	}
	return nil
}

// jsonb from Postgres
type JSONB map[string]interface{}

func (JSONB) GormDataType() string {
	return "jsonb"
}

func (j JSONB) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if j == nil {
		return clause.Expr{
			SQL:  "?",
			Vars: []interface{}{nil},
		}
	}

	valueString, err := json.Marshal(j)
	if err != nil {
		return clause.Expr{
			SQL:  "?",
			Vars: []interface{}{nil},
		}
	}

	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{string(valueString)},
	}
}

func (j *JSONB) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}

// func (j JSONB) MarshalJSON() ([]byte, error) {
// 	data, _ := interface{}(j).(map[string]interface{})

// 	res, err := json.Marshal(data)
// 	return res, err
// }

func (j *JSONB) UnmarshalJSON(b []byte) error {
	data := make(map[string]interface{}, 0)
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*j = JSONB(data)

	return nil
}

// NullableFloat для нулевых значений int
func NullableFloat(f float64) sql.NullFloat64 {
	if f != 0.0 {
		return sql.NullFloat64{Float64: f, Valid: true}
	}
	return sql.NullFloat64{}
}

// NullableInt export
func NullableInt(i int64) sql.NullInt64 {
	if i != 0 {
		return sql.NullInt64{Int64: i, Valid: true}
	}
	return sql.NullInt64{}
}

// NullableString export
func NullableString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
