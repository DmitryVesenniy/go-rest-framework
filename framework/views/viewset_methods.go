package views

type ViewSetMethod string

const (
	// Base methods
	GET  ViewSetMethod = "get"
	POST ViewSetMethod = "post"

	// Castom methods
	List     ViewSetMethod = "list"
	Create   ViewSetMethod = "create"
	Update   ViewSetMethod = "update"
	Retrieve ViewSetMethod = "retrieve"
	Delete   ViewSetMethod = "delete"
)
