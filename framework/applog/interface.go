package applog

type DiffInterface interface {
	ToDict() map[string]interface{}
	CalcDifference() error
}
