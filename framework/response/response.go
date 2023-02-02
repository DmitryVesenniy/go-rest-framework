package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DmitryVesenniy/go-rest-framework/framework/pagination"
	"github.com/DmitryVesenniy/go-rest-framework/framework/resterrors"
	"github.com/DmitryVesenniy/go-rest-framework/framework/serializers"
)

const (
	REST_DATA_KEY       = "data"
	REST_PAGINATION_KEY = "pagination"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func JsonResponceMap(w http.ResponseWriter, data map[string]interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		JsonErrorResponce(w, map[string][]string{resterrors.DEFAULT_KEY_ERROR: {err.Error()}}, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(b))
}

func JsonRespond(w http.ResponseWriter, data map[string]interface{}) {
	b, err := serializers.MarshallSerializer(data)
	if err != nil {
		JsonErrorResponce(w, map[string][]string{resterrors.DEFAULT_KEY_ERROR: {err.Error()}}, http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, string(b))
}

func JsonResponceMapSlice(w http.ResponseWriter, data map[string][]*string) {
	b, err := json.Marshal(data)
	if err != nil {
		JsonErrorResponce(w, map[string][]string{resterrors.DEFAULT_KEY_ERROR: {err.Error()}}, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(b))
}

func JsonErrorResponce(w http.ResponseWriter, errorList map[string][]string, status int) {
	w.WriteHeader(status)

	b, _ := json.Marshal(map[string]map[string][]string{"error": errorList})
	fmt.Fprintln(w, string(b))
}

func Response(w http.ResponseWriter, data interface{}) {
	resp := map[string]interface{}{
		REST_DATA_KEY: data,
	}
	JsonRespond(w, resp)
}

func ResponseList(w http.ResponseWriter, dataList interface{}, paginator *pagination.Pagination) {
	resp := map[string]interface{}{
		REST_DATA_KEY:       dataList,
		REST_PAGINATION_KEY: paginator,
	}

	JsonRespond(w, resp)
}
