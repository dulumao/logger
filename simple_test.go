package logger

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/ansel1/merry"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestFunc(t *testing.T) {
	DDStdout("哈哈", "hehe")
}

func TestSimple(t *testing.T) {
	const DEBUG = true

	log := &Logger{
		Hnd: new(SimpleLogHandler),
	}

	log.Hnd.SetConsole(DEBUG)

	if !DEBUG {
		log.Hnd.SetOutput(&lumberjack.Logger{
			Filename:  "logs/app.log",
			MaxSize:   1,
			MaxAge:    365,
			Compress:  true,
			LocalTime: true,
		})
	}

	//log.Hnd.SetOutput(os.Stderr)
	//logrus.Println("asda")

	for i := 0; i < 10000; i++ {
		log.Message(false, strconv.Itoa(i))
	}

	//log.Warn("http://%s", "127.0.0.1")
	//log.Info("http://%s", "127.0.0.1")
	//log.Debug("http://%s", "127.0.0.1")
	//
	//log.Error(testError(), "调用测试程序")

	defer func() {
		if err := recover(); err != nil {
			//log.Panic(fmt.Sprintf("%s\n", err))
			log.Error(err.(error), "出错了")
		}
	}()

	//if err := testError2(); err != nil {
	//	//frames := stack.Callers(3)
	//	panic(err)
	//}

	log.Warn("http://%s", "127.0.0.1")
	log.Info("http://%s", "127.0.0.1")
	log.Debug("http://%s", "127.0.0.1")
	log.Emer("Asdasd")
	log.Trace("Asdasd")

	log.Error(testError(), "调用测试程序")
}

func testError() error {
	return merry.Wrap(errors.New("出错了"))
	return merry.UserError("出错了")
}

func testError2() error {
	if err := recover(); err != nil {
		fmt.Printf("111: %s\n", err)
	}

	return merry.Wrap(errors.New("出错了"))
}

func TestWrap(t *testing.T) {
	log := &Logger{
		Hnd: Simple(),
	}

	log.Hnd.SetConsole(true)

	log.Print("测试内容")

	//log.Print(Wrap(errors.New("出错了")))
	//log.Print(WrapDebug(errors.New("出错了"), "1", "2"))
}
