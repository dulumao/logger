package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"runtime"
	"sync"

	"github.com/ansel1/merry"
	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
)

type LogLevel int

//等级	配置	释义	控制台颜色
//0	EMER	系统级紧急，比如磁盘出错，内存异常，网络不可用等	红色底
//1	ALRT	系统级警告，比如数据库访问异常，配置文件出错等	紫色
//2	CRIT	系统级危险，比如权限出错，访问异常等	蓝色
//3	EROR	用户级错误	红色
//4	WARN	用户级警告	黄色
//5	INFO	用户级重要	天蓝色
//6	DEBG	用户级调试	绿色
//7	TRAC	用户级基本输出	绿色

const (
	ERROR LogLevel = iota
	WARN
	INFO
	DEBUG
	PANIC
	FATAL
	EMER
	TRACE
)

type LogHandler interface {
	Log(LogLevel, error, string, ...interface{})
	Print(LogLevel, ...interface{})
	Message(bool, string)
	SetOutput(out io.Writer)
	SetConsole(color bool)
}

type SimpleLogHandler struct {
	isConsole bool
	logger    *log.Logger
	lock      sync.Mutex
}

func Simple() *SimpleLogHandler {
	var l = log.New(os.Stderr, "", log.LstdFlags)

	return &SimpleLogHandler{
		logger: l,
	}
}

func (h *SimpleLogHandler) Name(obj interface{}) string {
	return reflect.TypeOf(obj).Name()
}

func (h *SimpleLogHandler) SetConsole(isConsole bool) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.isConsole = isConsole
}

func (h *SimpleLogHandler) SetOutput(out io.Writer) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.logger.SetOutput(out)
}

func (h *SimpleLogHandler) StringifyLog(level LogLevel, err error, msg string, args ...interface{}) string {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	if err != nil {
		msg += ":\n" + merry.Details(err)
	}

	return msg
}

func (h *SimpleLogHandler) AddLevelPrevix(level LogLevel, text string) string {
	switch level {
	case ERROR:
		text = "[ERROR] " + text
	case WARN:
		text = "[WARN] " + text
	case INFO:
		text = "[INFO] " + text
	case DEBUG:
		text = "[DEBUG] " + text
	case PANIC:
		text = "[PANIC] " + text
	case FATAL:
		text = "[FATAL] " + text
	case EMER:
		text = "[EMER] " + text
	case TRACE:
		text = "[TRACE] " + text
	}
	return text
}

func (h *SimpleLogHandler) AddLevelColor(level LogLevel, text string) string {
	switch level {
	case DEBUG:
		return "\033[90m" + text + "\033[0m"
	case WARN:
		return "\033[93m" + text + "\033[0m"
	case INFO:
		return "\033[1;34m" + text + "\033[0m"
	case ERROR:
		return "\033[91m" + text + "\033[0m"
	case EMER:
		return "\033[1;36m" + text + "\033[0m"
	case FATAL:
		return "\033[1;35m" + text + "\033[0m"
	case PANIC:
		return "\033[1;41m" + text + "\033[0m"
	case TRACE:
		return "\033[1;32m" + text + "\033[0m"
	}

	return text
}

func (h *SimpleLogHandler) StringifyMessage(isIncoming bool, msg string) string {
	var text string

	if isIncoming {
		text = ">>> " + msg
	} else {
		text = "<<< " + msg
	}

	return text
}

func (h *SimpleLogHandler) Log(level LogLevel, err error, msg string, args ...interface{}) {
	h.lock.Lock()
	defer h.lock.Unlock()

	text := h.StringifyLog(level, err, msg, args...)
	text = h.AddLevelPrevix(level, text)

	if h.isConsole {
		text = h.AddLevelColor(level, text)
	}

	h.logger.Print(text)
}

func (h *SimpleLogHandler) Print(level LogLevel, msg ...interface{}) {
	h.lock.Lock()
	defer h.lock.Unlock()

	text := h.AddLevelPrevix(level, fmt.Sprint(msg...))

	if h.isConsole {
		text = h.AddLevelColor(level, text)
	}

	h.logger.Print(text)
}

func (h *SimpleLogHandler) Message(isIncoming bool, msg string) {
	h.Log(DEBUG, nil, h.StringifyMessage(isIncoming, msg))
}

type Logger struct {
	Hnd LogHandler
}

func (l Logger) Error(err error, msg string, args ...interface{}) {
	l.Hnd.Log(ERROR, err, msg, args...)
}

func (l Logger) Errorf(msg string, args ...interface{}) {
	l.Hnd.Log(ERROR, nil, msg, args...)
}

func (l Logger) Warn(msg string, args ...interface{}) {
	l.Hnd.Log(WARN, nil, msg, args...)
}

func (l Logger) Info(msg string, args ...interface{}) {
	l.Hnd.Log(INFO, nil, msg, args...)
}

func (l Logger) Debug(msg string, args ...interface{}) {
	l.Hnd.Log(DEBUG, nil, msg, args...)
}

func (l Logger) Panic(err error) {
	l.Hnd.Log(PANIC, nil, err.Error())
}

func (l Logger) Paninf(msg string, args ...interface{}) {
	l.Hnd.Log(PANIC, nil, msg, args...)
}

func (l Logger) Fatal(err error) {
	l.Hnd.Log(FATAL, nil, err.Error())
	os.Exit(1)
}

func (l Logger) Emer(msg string, args ...interface{}) {
	l.Hnd.Log(EMER, nil, msg, args...)
}

func (l Logger) Trace(msg string, args ...interface{}) {
	l.Hnd.Log(TRACE, nil, msg, args...)
}

func (l Logger) Printf(msg string, args ...interface{}) {
	l.Hnd.Log(DEBUG, nil, msg, args...)
}

func (l Logger) Print(args ...interface{}) {
	l.Hnd.Print(DEBUG, args...)
}

func (l Logger) Message(isIncoming bool, message string) {
	l.Hnd.Message(isIncoming, message)
}

func Wrap(err error) error {
	pc, filename, linenr, _ := runtime.Caller(1)

	return merry.Here(err).WithMessagef("\n\nerror in function[%s] file[%s] line[%d]", runtime.FuncForPC(pc).Name(), filename, linenr)
}

func WrapDebug(err error, vars ...interface{}) error {
	pc, filename, linenr, _ := runtime.Caller(1)

	return merry.Here(err).WithMessagef(`
error in function[%s] file[%s] line[%d]
↓↓↓ Debug variables ↓↓↓
--------------------------------------------------------------------------------
%s--------------------------------------------------------------------------------
cause of error`, runtime.FuncForPC(pc).Name(), filename, linenr, spew.Sdump(vars...))
}

func DDStdout(vars ...interface{}) {
	pc, filename, linenr, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stdout, `
function[%s] file[%s] line[%d]
↓↓↓ Debug variables ↓↓↓
--------------------------------------------------------------------------------
%s--------------------------------------------------------------------------------
`, runtime.FuncForPC(pc).Name(), filename, linenr, spew.Sdump(vars...))
}

// JSStdout json dump
func JSStdout(i ...interface{}) {
	pc, filename, linenr, _ := runtime.Caller(1)
	formatter := prettyjson.NewFormatter()

	fmt.Fprintf(os.Stdout, `
function[%s] file[%s] line[%d]
↓↓↓ Debug variables ↓↓↓
--------------------------------------------------------------------------------
`, runtime.FuncForPC(pc).Name(), filename, linenr)

	formatter.KeyColor = color.New(color.FgWhite)
	formatter.StringColor = color.New(color.FgGreen)
	formatter.BoolColor = color.New(color.FgYellow)
	formatter.NumberColor = color.New(color.FgCyan)
	formatter.NullColor = color.New(color.FgMagenta)

	formatter.Indent = 4
	formatter.DisabledColor = false

	for _, v := range i {
		if s, err := formatter.Marshal(v); err == nil {
			fmt.Fprintln(os.Stdout, string(s))
		} else {
			fmt.Fprintln(os.Stdout, "%s", err.Error())
		}
		//buffer := &bytes.Buffer{}
		//encoder := json.NewEncoder(buffer)
		//encoder.SetEscapeHTML(false)
		//encoder.SetIndent("", "\t")
		//if err := encoder.Encode(v); err == nil {
		//	fmt.Print(buffer.String())
		//} else {
		//	fmt.Errorf("%s", err.Error())
		//}
	}

	fmt.Fprintln(os.Stdout, `--------------------------------------------------------------------------------`)
}

func Dump(w io.Writer, err error) {
	fmt.Fprintf(w, err.Error())
}

func DumpVars(w io.Writer, vars ...interface{}) {
	pc, filename, linenr, _ := runtime.Caller(1)
	fmt.Fprintf(w, `
function[%s] file[%s] line[%d]
↓↓↓ Debug variables ↓↓↓
--------------------------------------------------------------------------------
%s--------------------------------------------------------------------------------
`, runtime.FuncForPC(pc).Name(), filename, linenr, spew.Sdump(vars...))
}

func WrapReturn(err error) func() error {
	return func() error {
		if err != nil {
			pc, filename, linenr, _ := runtime.Caller(0)

			return merry.Here(err).WithMessagef("\n\nerror in function[%s] file[%s] line[%d]", runtime.FuncForPC(pc).Name(), filename, linenr)
		}

		return nil
	}
}

// WrapReturnMulti lorem ipsum
func WrapReturnMulti(err *error) func() {
	return func() {
		if *err != nil {
			pc, filename, linenr, _ := runtime.Caller(0)

			*err = merry.Here(*err).WithMessagef("\n\nerror in function[%s] file[%s] line[%d]", runtime.FuncForPC(pc).Name(), filename, linenr)

			return
		}

		return
	}
}
