package logger

import "sync"

var single *Logger
var once sync.Once

func Console() *Logger {
	once.Do(func() {
		logger := &Logger{
			Hnd: Simple(),
		}

		logger.Hnd.SetConsole(true)

		single = logger
	})

	return single
}
