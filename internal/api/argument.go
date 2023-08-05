package api

type Argument interface {
	Name() string
	Description() string
	Extract(map[string]string, []string) []string
	String(int) string
}
