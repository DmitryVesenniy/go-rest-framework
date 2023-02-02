package resterrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	AppErr              = &AppError{Err: "Server error"}
	NotFoundErr         = &NotFoundError{}
	PermissionDeniedErr = &PermissionDeniedError{Err: "Permission denied"}
	ValidErr            = &ValidError{}
	InvalidPasswordErr  = &InvalidPasswordError{}
)

/*====== NotFoundError ======*/
type BaseError struct {
	Err string
}

func (b *BaseError) Error() string {
	if b.Err != "" {
		return b.Err
	}
	return "Bad response"
}
func (BaseError) Status() int {
	return http.StatusBadRequest
}
func (e *BaseError) RestError() map[string][]string {
	return map[string][]string{DEFAULT_KEY_ERROR: {e.Error()}}
}

/*====== PermissionDeniedError ======*/
type AppError struct {
	BaseError
	Err string
}

func (e *AppError) Error() string {
	return e.Err
}
func (AppError) Status() int {
	return http.StatusInternalServerError
}
func (e *AppError) RestError() map[string][]string {
	return map[string][]string{DEFAULT_KEY_ERROR: {e.Error()}}
}

/*====== NotFoundError ======*/
type NotFoundError struct{ BaseError }

func (NotFoundError) Error() string {
	return "Not found"
}
func (NotFoundError) Status() int {
	return http.StatusNotFound
}
func (e *NotFoundError) RestError() map[string][]string {
	return map[string][]string{DEFAULT_KEY_ERROR: {e.Error()}}
}

/*====== PermissionDeniedError ======*/
type PermissionDeniedError struct {
	BaseError
	Err string
}

func (p *PermissionDeniedError) Error() string {
	return p.Err
}
func (PermissionDeniedError) Status() int {
	return http.StatusForbidden
}
func (e *PermissionDeniedError) RestError() map[string][]string {
	return map[string][]string{DEFAULT_KEY_ERROR: {e.Error()}}
}

/*====== PermissionDeniedError ======*/
type ValidError struct {
	BaseError
	Field string
}

func (v *ValidError) Error() string {
	return fmt.Sprintf("Invalid field '%s'", v.Field)
}
func (ValidError) Status() int {
	return http.StatusUnprocessableEntity
}
func (e *ValidError) RestError() map[string][]string {
	return map[string][]string{DEFAULT_KEY_ERROR: {e.Error()}}
}

/*====== ModelInstanceError ======*/
type ModelInstanceError struct {
	BaseError
	Err error
}

func (v *ModelInstanceError) Error() string {
	return fmt.Sprintf("ModelInstanceError: '%s'", v.Err.Error())
}
func (ModelInstanceError) Status() int {
	return http.StatusUnprocessableEntity
}
func (e *ModelInstanceError) RestError() map[string][]string {
	return map[string][]string{DEFAULT_KEY_ERROR: {e.Error()}}
}

/*====== PermissionDeniedError ======*/
type InvalidPasswordError struct {
	BaseError
}

func (v *InvalidPasswordError) Error() string {
	return "Invalid password"
}
func (InvalidPasswordError) Status() int {
	return http.StatusUnprocessableEntity
}
func (e *InvalidPasswordError) RestError() map[string][]string {
	return map[string][]string{DEFAULT_KEY_ERROR: {e.Error()}}
}

func RestErrorResponce(w http.ResponseWriter, restError RestError) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(restError.Status())
	b, _ := json.Marshal(map[string]map[string][]string{"error": restError.RestError()})
	fmt.Fprintln(w, string(b))
}

/*====== UnauthorizedError ======*/
type UnauthorizedError struct {
	BaseError
	Err string
}

func (p *UnauthorizedError) Error() string {
	return p.Err
}
func (UnauthorizedError) Status() int {
	return http.StatusUnauthorized
}
func (e *UnauthorizedError) RestError() map[string][]string {
	return map[string][]string{DEFAULT_KEY_ERROR: {e.Error()}}
}
