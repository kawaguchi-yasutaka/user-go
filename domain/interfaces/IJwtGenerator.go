package interfaces

type IJwtGenerator interface {
	GenerateJwtToken(payload map[string]interface{}) (string, error)
}
