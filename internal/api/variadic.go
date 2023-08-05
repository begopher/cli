package api

type Variadic interface {
	Arg() string
	Allowed() bool
	Extract([]string) ([]string, error)
	String() string
}
