package constants

// Metrics constants
var Metrics = struct {
	APIRequestCount string
	APIRequestPath  string
}{
	APIRequestCount: "api_request_count",
	APIRequestPath:  "api_request_path",
}

// AppInfo constants
var AppInfo = struct {
	Name    string
	Version string
}{
	Name:    "Will.IAM",
	Version: "2.0",
}
