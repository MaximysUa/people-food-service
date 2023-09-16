package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
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

func GetLogger() *Logger {
	return &Logger{e}
}
func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}

}

// Создаём новую сущность логгирования и настраиваем её
func Init() {
	l := logrus.New()
	l.SetReportCaller(true)
	//формат возвращаемого значения - текст, так же может быть и jsong
	l.Formatter = &logrus.TextFormatter{
		//получаем фрейм в котором происходит логирование, в нем есть информация о фаиле в котором происходит логирование
		//в строчке -> там есть инфа о линии и функц в в которой что то происходит
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", fileName, frame.Line)
		},
		//отключаем цвета (зачем?) TODO поиграться с цветами
		DisableColors: true,
		FullTimestamp: true,
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

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
