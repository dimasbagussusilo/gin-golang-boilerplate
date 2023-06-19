package forms

type HealthCheckResponse struct {
	ServerStatus   string `json:"server_status"`
	DatabaseStatus string `json:"database_status"`
	DatabaseName   string `json:"database_name"`
	DatabaseHost   string `json:"database_host"`
}
