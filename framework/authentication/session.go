package authentication

import (
	"net/http"

	"github.com/DmitryVesenniy/go-rest-framework/framework/utils"
)

type SessionInterface interface {
	GenerateKey() (string, error)
	GetSession(r *http.Request) map[string]interface{}
	GetValue(key string, r *http.Request) interface{}
	SetValue(key string, v interface{}, r *http.Request)
}

type SessionBase struct {
}

func (s *SessionBase) GenerateKey() (string, error) {
	return utils.RandomGenerator(utils.SESSION_LENGTH_KEY, utils.SESSION_PREFIX), nil
}
func (s *SessionBase) GetSession(r *http.Request) map[string]interface{} {
	return map[string]interface{}{}
}
func (s *SessionBase) GetValue(key string, r *http.Request) interface{} {
	return nil
}
func (s *SessionBase) SetValue(key string, v interface{}, r *http.Request) {

}
