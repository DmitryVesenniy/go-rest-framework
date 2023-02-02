package authentication

import "net/http"

type AccessLevel uint

const (
	ADMIN  AccessLevel = 1
	WRITER AccessLevel = 2
	READER AccessLevel = 3
)

var (
	HTTP_METHOD_ACCESS = map[string]AccessLevel{
		http.MethodGet:     READER,
		http.MethodHead:    READER,
		http.MethodOptions: READER,
		http.MethodPost:    WRITER,
		http.MethodPut:     WRITER,
		http.MethodPatch:   WRITER,
		http.MethodDelete:  ADMIN,
	}
)

func IsAllow(method string, access AccessLevel) bool {
	minimalAccessMethod, ok := HTTP_METHOD_ACCESS[method]
	if !ok {
		return false
	}

	return minimalAccessMethod >= access
}
