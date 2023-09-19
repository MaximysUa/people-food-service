package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"people-food-service/iternal/config"
	"runtime"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// сущность для записи сразу в файл и оутпут
type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

// метод будет вызываться каждый раз для записи
func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

// будет возвращать левелы
func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// Код на случай если нам необходимо будет создать еще 1 сущность логгера
// по умолчанию логгер - синглтон
var e *logrus.Entry

// удобно тем, что благодаря этой структуре можно безболезненно заменить логрус на любой другой логгер
type Logger struct {
	*logrus.Entry
}

// GetLogger We use it after a method init to get a configurated entity of logger
func GetLogger() *Logger {
	return &Logger{e}
}
func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}

}

// Init and configure the entity of logrus
func Init(cfg *config.Config) {
	l := logrus.New()
	l.SetReportCaller(true)
	//формат возвращаемого значения - текст, так же может быть и json
	l.Formatter = &logrus.TextFormatter{
		//получаем фрейм в котором происходит логирование, в нем есть информация о фаиле в котором происходит логирование
		//в строчке -> там есть инфа о линии и функц в в которой что то происходит
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", fileName, frame.Line)
		},
		//TODO поиграться с цветами
		ForceColors:               true,
		FullTimestamp:             true,
		EnvironmentOverrideColors: true,
	}

	// создаём папку для хранения логов
	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}
	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	//это значит: ничего никуда не пиши
	l.SetOutput(io.Discard)
	//создаём крюки для записи в разные места
	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	level, err := logrus.ParseLevel(cfg.Env)
	if err != nil {
		log.Fatal(err)
	}
	logrus.SetLevel(level)
	log.Printf("Logger level is %s", level)
	e = logrus.NewEntry(l)
}
