package envexample

import (
	"os"
)

// SetEnv function
func SetEnv() {
	os.Setenv("DB_CONNECTION", "mysql")
	os.Setenv("DB_HOST", "db")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_DATABASE", "mysql")
	os.Setenv("DB_USERNAME", "root")
	os.Setenv("DB_PASSWORD", "root")

	os.Setenv("AWS_S3_REGION", "sa-east-1")
	os.Setenv("AWS_S3_BUCKET", "example-bucket")

	os.Setenv("AWS_SES_SENDER", "example@example.com")
}
