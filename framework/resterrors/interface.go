package resterrors

type RestError interface {
	Error() string
	Status() int
	RestError() map[string][]string
}
