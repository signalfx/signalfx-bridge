package metrics

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/caarlos0/env"
)

type Config struct {
	CloudFoundryApiURL    string `env:"CLOUDFOUNDRY_API_URL,required"`
	CFUAAURL              string `env:"CF_UAA_URL,required"`
	CFUsername            string `env:"CF_USERNAME,required"`
	CFPassword            string `env:"CF_PASSWORD,required"`
	InsecureSSLSkipVerify bool   `env:"INSECURE_SSL_SKIP_VERIFY" envDefault:"false"`
	EnableTSDBServer      bool   `env:"ENABLE_TSDB_SERVER" envDefault:"true"`

	BoshDirectorURL string `env:"BOSH_DIRECTOR_URL,required"`
	BoshUsername    string `env:"BOSH_CLIENT_ID,required"`
	BoshPassword    string `env:"BOSH_CLIENT_SECRET,required"`

	// This will be populated automatically in the main package if not supplied
	TrafficControllerURL          string   `env:"TRAFFIC_CONTROLLER_URL" envDefault:""`
	FirehoseSubscriptionID        string   `env:"FIREHOSE_SUBSCRIPTION_ID" envDefault:"signalfx"`
	FlushIntervalSeconds          int      `env:"FLUSH_INTERVAL_SECONDS" envDefault:"3"`
	FirehoseIdleTimeoutSeconds    int      `env:"FIREHOSE_IDLE_TIMEOUT_SECONDS" envDefault:"20"`
	FirehoseReconnectDelaySeconds int      `env:"FIREHOSE_RECONNECT_DELAY_SECONDS" envDefault:"5"`
	DeploymentsToInclude          []string `env:"DEPLOYMENTS_TO_INCLUDE" envDefault:"" envSeparator:";"`
	MetricsToExclude              []string `env:"METRICS_TO_EXCLUDE" envDefault:"" envSeparator:";"`

	AppMetadataCacheExpirySeconds int `env:"APP_METADATA_CACHE_EXPIRY_SECONDS" envDefault:"300"`

	SignalFxIngestURL   string `env:"SIGNALFX_INGEST_URL"`
	SignalFxAccessToken string `env:"SIGNALFX_ACCESS_TOKEN,required"`
}

func GetConfigFromEnv() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)

	for i, v := range cfg.DeploymentsToInclude {
		cfg.DeploymentsToInclude[i] = strings.TrimSpace(v)
	}

	for i, v := range cfg.MetricsToExclude {
		cfg.MetricsToExclude[i] = strings.TrimSpace(v)
	}

	return &cfg, err
}

func (cfg *Config) ScrubbedString() string {
	v := reflect.ValueOf(*cfg)

	values := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		if typeField.Name == "CFPassword" ||
			typeField.Name == "BoshPassword" {
			values[typeField.Name] = "***********"
		} else {
			values[typeField.Name] = v.Field(i).Interface()
		}
	}

	return fmt.Sprintf("%#v", values)
}
