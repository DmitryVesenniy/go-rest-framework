package serializers

const (
	SerializerKeyRequired  = "required"
	SerializerKeyReadOnly  = "read_only"
	SerializerKeyWriteOnly = "write_only"
	SerializerKeyAllowNull = "allow_null"
)

type Rules struct {
	Required  bool
	ReadOnly  bool
	WriteOnly bool
	AllowNull bool
	Validator []func(interface{}) bool
}
