package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type LogConfig struct {
	Level  string `envconfig:"LOG_LEVEL" default:"info"`
	Format string `envconfig:"LOG_FORMAT" default:"json"`
}

type Config struct {
	Log     LogConfig
	DataDog DataDogConfig
	AWS     AWSConfig
	App     AppConfig
}

type AWSConfig struct {
	SQSTasksQueue         string        `envconfig:"SQS_TASKS" required:"true"`
	SQSResultsQueue       string        `envconfig:"SQS_RESULTS"` // needed only for testing
	SNSResultsTopic       string        `envconfig:"SNS_RESULTS" required:"true"`
	S3BucketName          string        `envconfig:"S3_BUCKET_FILES" required:"true"`
	SQSPollingNumMessages int           `envconfig:"SQS_POLLING_NUM_MESSAGES" default:"5"`
	SQSPollingTimeout     time.Duration `envconfig:"SQS_POLLING_TIMEOUT" default:"10s"`
}

type AppConfig struct {
	PDFToolsLicenseKey string        `envconfig:"PDFTOOLS_LICENSE_KEY" required:"true"`
	Environment        string        `envconfig:"ENVIRONMENT" required:"true"`
	ProcessingTimeout  time.Duration `envconfig:"PROCESSING_TIMEOUT" default:"10m"`
	ConcurrentTasks    int           `envconfig:"CONCURRENT_TASKS" default:"1"`
	SendResultEnabled  bool          `envconfig:"SEND_RESULT_ENABLED" default:"true"`
}
type DataDogConfig struct {
	InstrumentationEnabled bool `envconfig:"DD_INSTRUMENTATION_ENABLED" default:"true"`
}

func Load(prefix string) (Config, error) {
	c := Config{}

	err := envconfig.Process(prefix, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
