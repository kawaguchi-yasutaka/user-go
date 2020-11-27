package interfaces

type IJwtGenerator interface {
	Generate(payload map[string]interface{}) (string, error)
}
