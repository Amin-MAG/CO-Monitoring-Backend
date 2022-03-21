package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	. "os"
)

// Time format for our log package
const iso8601DateFormat = "2006-01-02T15:04:05"

// Config
// It specifies the details about a new Logger
type Config struct {
	IsProductionMode bool `default:"true"`
	IsReportCaller   bool `default:"false"`
	OutputFileConfig *OutputFileConfig
}

// OutputFileConfig
// To specify more details about saving logs in files
type OutputFileConfig struct {
	Path string
	Name string
	// Default is storing logs in 100KB files
	MaxSizeInBytes int64
	// Default is storing logs in 10 files
	MaxNumberOfFiles int `default:"10"`
}

func (ofc OutputFileConfig) FullPath(fileIndex int) string {
	if ofc.Name == "" {
		ofc.Name = "log_file"
	}

	// Calculate number of digits
	var numberOfDigits int
	if ofc.MaxNumberOfFiles == 1 {
		numberOfDigits = 1
	} else {
		numberOfDigits = int(math.Log10(float64(ofc.MaxNumberOfFiles)-1) + 1)
	}

	return fmt.Sprintf("%s/%s-%0*d.log", ofc.Path, ofc.Name, numberOfDigits, fileIndex)
}

// Logger
// It is a wrapper for logrus.Logger.
type Logger struct {
	*logrus.Logger
	config Config
}

// NewLogger
// It is to create new Logger by a Config.
func NewLogger(config Config) (*Logger, error) {
	// Create new Logrus logger
	ll := logrus.New()
	// TODO: Make this log level env
	ll.SetLevel(logrus.DebugLevel)
	ll.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: iso8601DateFormat,
	})
	// Adds file and line number
	ll.SetReportCaller(config.IsReportCaller)
	// Check if it is going to be saved on file
	if config.OutputFileConfig != nil {
		f, err := OpenFile((*config.OutputFileConfig).FullPath(0), O_RDWR|O_CREATE|O_APPEND, 0666)
		if err != nil {
			fmt.Println("Could Not Open Log File : " + err.Error())
		}
		ll.SetOutput(f)

		// For checking storage limitation
		ll.AddHook(NewFileThresholdHook(*config.OutputFileConfig))
	}
	ll.Debugf("Ther logger has been created: %+v\n", config)

	return &Logger{
		ll,
		config,
	}, nil
}
