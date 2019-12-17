package logging

import (
	"log"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ConfigLogging definie configuração inicial do Appender
type ConfigLogging interface {
	Configure()
}

// FileAppender configura logging para escrita em arquivo
type FileAppender struct {
	config config.Configuration
}

// NewFileAppender cria FileAppender
func NewFileAppender(config config.Configuration) (cl FileAppender) {
	cl = FileAppender{
		config: config,
	}
	return
}

// Configure define configurações de escrita em arquivo
func (cl FileAppender) Configure() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   cl.config.Logging.File,
		MaxSize:    1,
		MaxBackups: 14,
		MaxAge:     28,
	})
}
