package authentication

type ErrorNotFoundRole struct {
}

func (e *ErrorNotFoundRole) Error() string {
	return "Not found role"
}

type ErrorNotFoundOrganisation struct {
}

func (e *ErrorNotFoundOrganisation) Error() string {
	return "Not found organisation"
}
