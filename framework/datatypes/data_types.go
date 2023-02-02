package datatypes

type DictInterface interface {
	ToDict() map[string]interface{}
	Set(key string, value interface{})
	Get(key string) interface{}
}
