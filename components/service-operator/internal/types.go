package internal

type Action string

const (
	Create Action = "CREATE"
	Update Action = "UPDATE"
	Delete Action = "DELETE"
	Retry  Action = "RETRY"
)

type BasicAuth struct {
	Username string
	Password string
}
