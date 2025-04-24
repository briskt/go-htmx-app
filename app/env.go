package app

import (
	_ "embed"

	"github.com/kelseyhightower/envconfig"
)

//go:embed commit.txt
var Commit string

var Env struct {
	// AppEnv is used to help switch settings based on where the application is being run.
	// Set the default to "prod" for safety in case it is not set correctly.
	AppEnv         string   `split_words:"true" default:"prod"`
	AppName        string   `split_words:"true" default:"Go HTMX"`
	AppURL         string   `split_words:"true" default:"http://localhost:8100"`
	APIAccessKeys  []string `split_words:"true"`
	BrandColor     string   `split_words:"true" default:"#f57c00"`
	DisableTLS     bool     `split_words:"true"`
	EmailService   string   `split_words:"true" default:"fake"`
	EmailSignature string   `split_words:"true" default:"This was sent by an automated process. Please do not reply."`
	FromEmail      string   `split_words:"true" default:"no_reply@example.com"`
	HelpCenterURL  string   `split_words:"true" default:"https://example.com"`
	LogLevel       string   `split_words:"true" default:"debug"`
	SandboxEmail   string   `split_words:"true" default:""`
	SupportEmail   string   `split_words:"true" default:"support@example.com"`
	SupportName    string   `split_words:"true" default:"Help Desk"`

	SessionSecret string `split_words:"true"`

	AWSAccessKeyID     string `split_words:"true"`
	AWSRegion          string `split_words:"true"`
	AWSSecretAccessKey string `split_words:"true"`

	MailgunDomain string `split_words:"true"`
	MailgunAPIKey string `split_words:"true"`

	PostgresUser     string `split_words:"true"`
	PostgresPassword string `split_words:"true"`
	PostgresHost     string `split_words:"true"`
	PostgresDB       string `split_words:"true"`

	SamlSpEntityID                  string `split_words:"true"`
	SamlSpCert                      string `split_words:"true"`
	SamlSpPrivateKey                string `split_words:"true"`
	SamlAssertionConsumerServiceURL string `split_words:"true"`
	SamlIdpMetadataURL              string `split_words:"true"`
}

// readEnv loads environment data into `Env`
func readEnv() {
	err := envconfig.Process("", &Env)
	if err != nil {
		panic("error loading env vars: " + err.Error())
	}
}
