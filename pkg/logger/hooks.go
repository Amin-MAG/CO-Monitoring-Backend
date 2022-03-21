package logger

import (
	"github.com/sirupsen/logrus"
	. "os"
)

// FileThresholdHook
// It's a hook to make the storage usage for logging stable
type FileThresholdHook struct {
	logrus.Hook
	OutputFileConfig OutputFileConfig
}

// NewFileThresholdHook
// To create a new hook
func NewFileThresholdHook(ofc OutputFileConfig) FileThresholdHook {
	return FileThresholdHook{
		OutputFileConfig: ofc,
	}
}

// Fire
// Implements the fire interface of hook
func (hook FileThresholdHook) Fire(entry *logrus.Entry) (err error) {
	_ = hook.checkThreshold(hook.OutputFileConfig, entry)

	return nil
}

// Levels
// Implements the levels interface of hook
func (hook FileThresholdHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.TraceLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (hook FileThresholdHook) shiftAllLogFiles(config OutputFileConfig) error {
	// Shift all log files
	for i := config.MaxNumberOfFiles - 1; i > 0; i-- {
		_ = Rename(config.FullPath(i-1), config.FullPath(i))
	}

	// Clean the first log file
	f, err := OpenFile(config.FullPath(0), O_RDWR|O_CREATE|O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

// checkThreshold
// Check the first log file and will shift the files if this
// main logging file excess the threshold.
func (hook FileThresholdHook) checkThreshold(config OutputFileConfig, entry *logrus.Entry) error {
	// Get file size
	f, err := Stat(config.FullPath(0))
	if err != nil {
		return err
	}

	// Check threshold
	if f.Size() > config.MaxSizeInBytes {
		_ = hook.shiftAllLogFiles(config)
		mainLogFile, _ := OpenFile(config.FullPath(0), O_RDWR|O_CREATE|O_APPEND, 0666)
		entry.Logger.SetOutput(mainLogFile)
	}

	return nil
}
