package api

type Arguments interface {
	Names() []string
	Extract(map[string]string, []string) ([]string, error)
	Count() int
	String() string
}
