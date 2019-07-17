// logs
package logs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	defaultLoggerFlag  = "default"
	OsStdOutLoggerType = iota
	FileLoggerType

	Ldate = 1 << iota
	Ltime
	Lmicroseconds
	Llongfile
	Lshortfile
	LUTC
	LstdFlags = Ldate | Ltime

	defaultFileLoggerFlag     = "snow"
	defaultFileLoggerPath     = "./"
	defaultFileLoggerTimeFlag = "20060102"
	defaultLoggerHeader       = LstdFlags

	TimeStampLoggerHeader = LstdFlags
	FileNameLoggerHeader  = Lshortfile

	debugFlag    = "[DEBUG] "
	infoFlag     = "[INFO] "
	warnFlag     = "[WARN] "
	errorFlag    = "[ERROR] "
	criticalFlag = "[CRIT] "

	fatalFlag = "[FATAL] "
	panicFlag = "[PANIC] "
)

var defaultCallDepth = 3
var loggerPool *map[string]*Logger
var loggers map[string]*Logger
var defaultLogger *Logger
var defaultLoggerRegister = StdOutLoggerRegister{
	loggerRegister{
		lType:      OsStdOutLoggerType,
		flag:       defaultLoggerFlag,
		headerFlag: defaultLoggerHeader,
	},
}

type Logger struct {
	writeCloser        io.WriteCloser
	getWriteCloserFunc func() (io.WriteCloser, error)
	rwLock             sync.RWMutex
	lock               sync.Mutex
	register           registerInterface
	buf                []byte
}

type loggerRegister struct {
	lType      int
	flag       string
	headerFlag int
}

type FileLoggerRegister struct {
	loggerRegister
	timeFlag string
	path     string
}

type StdOutLoggerRegister struct {
	loggerRegister
}

type registerInterface interface {
	register() (*Logger, error)
	getLType() int
	getFlag() string
	getHeaderFlag() int
	setHeaderFlag(int)
}

func init() {

	loggers = make(map[string]*Logger)
	loggerPool = &loggers

	if defaultLogger == nil {
		setDefaultLogger(&defaultLoggerRegister)
	}
	go startLoggerRollingWatcher()

}
func (this *loggerRegister) getLType() int {
	return this.lType
}

func (this *loggerRegister) getFlag() string {
	return this.flag
}

func (this *loggerRegister) getHeaderFlag() int {
	return this.headerFlag
}

func (this *loggerRegister) setHeaderFlag(headerFlag int) {
	this.headerFlag = headerFlag
}

func NewFileLoggerRegister(lType int, flag string, headerFlag int, timeFlag string, path string) (loggerRegiser *FileLoggerRegister) {
	loggerRegiser = new(FileLoggerRegister)
	loggerRegiser.loggerRegister.lType = lType
	loggerRegiser.loggerRegister.flag = flag
	loggerRegiser.loggerRegister.headerFlag = headerFlag
	loggerRegiser.timeFlag = timeFlag
	loggerRegiser.path = path
	return
}

func NewStdOutLoggerRegister(lType int, flag string, headerFlag int) (loggerRegister *StdOutLoggerRegister) {
	loggerRegister = new(StdOutLoggerRegister)
	loggerRegister.loggerRegister.lType = lType
	loggerRegister.loggerRegister.flag = flag
	loggerRegister.loggerRegister.headerFlag = headerFlag
	return

}

func (this *FileLoggerRegister) register() (logger *Logger, err error) {

	logger = new(Logger)

	if this.path == "" {
		this.path = defaultFileLoggerPath
	}

	if this.flag == "" {
		this.flag = defaultFileLoggerFlag
	}

	if this.timeFlag == "" {
		this.timeFlag = defaultFileLoggerTimeFlag
	}

	logger.getWriteCloserFunc = func() (io.WriteCloser, error) {
		logNameTS := time.Now().Format(this.timeFlag)
		return os.OpenFile(
			filepath.Join(this.path, this.flag+"."+logNameTS+".log"),
			os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	}

	logger.writeCloser, err = logger.getWriteCloserFunc()
	logger.register = this

	return
}

func (this *StdOutLoggerRegister) register() (logger *Logger, err error) {
	logger = new(Logger)
	logger.getWriteCloserFunc = func() (io.WriteCloser, error) {
		return os.Stdout, nil
	}
	logger.writeCloser, err = logger.getWriteCloserFunc()
	logger.register = this
	return
}

func SetDefaultCallDepth(depth int) {
	defaultCallDepth = depth
}

func SetDefaultHeaderFLag(headerFlag int) {
	defaultLogger.register.setHeaderFlag(headerFlag)
}

func SetDefaultLogger(register registerInterface) (err error) {
	return setDefaultLogger(register)
}

func setDefaultLogger(register registerInterface) (err error) {
	defaultLogger, err = register.register()
	if err != nil {
		return
	}

	(*loggerPool)[register.getFlag()] = defaultLogger

	return
}

func regLoggerMulti(registers ...registerInterface) (err error) {

	for _, register := range registers {
		logger, err := register.register()
		if err != nil {
			return err
		}
		(*loggerPool)[register.getFlag()] = logger

	}

	return
}

func Init(registers ...registerInterface) (err error) {
	err = regLoggerMulti(registers...)
	return
}

func startLoggerRollingWatcher() {
	for {
		for _, logger := range loggers {
			if logger.register.getLType() == FileLoggerType {
				register := logger.register.(*FileLoggerRegister)
				logNameTS := time.Now().Format(register.timeFlag)
				fileName := filepath.Join(register.path, register.flag+"."+logNameTS+".log")
				file := logger.writeCloser.(*os.File)
				if fileName != file.Name() {
					logger.rwLock.Lock()
					logger.writeCloser.Close()
					logger.writeCloser, _ = logger.getWriteCloserFunc()
					logger.rwLock.Unlock()
				} else {
					continue
				}
			} else {
				continue
			}
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (this *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	headerFlag := this.register.getHeaderFlag()
	if headerFlag&LUTC != 0 {
		t = t.UTC()
	}
	if headerFlag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if headerFlag&Ldate != 0 {
			year, month, day := t.Date()
			*buf = append(*buf, '[')
			itoa(buf, year, 4)
			*buf = append(*buf, '-')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '-')
			itoa(buf, day, 2)
			*buf = append(*buf, ']')
			*buf = append(*buf, ' ')
		}
		if headerFlag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			*buf = append(*buf, '[')
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if headerFlag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ']')
			*buf = append(*buf, ' ')
		}
	}
	if headerFlag&(Lshortfile|Llongfile) != 0 {
		if headerFlag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func (this *Logger) output(calldepth int, level string, s string) error {
	headerFlag := this.register.getHeaderFlag()
	now := time.Now() // get this early.
	var file string
	var line int
	this.lock.Lock()
	defer this.lock.Unlock()
	if headerFlag&(Lshortfile|Llongfile) != 0 {
		// release lock while getting caller info - it's expensive.
		this.lock.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		this.lock.Lock()
	}
	this.buf = this.buf[:0]
	this.formatHeader(&this.buf, now, file, line)
	this.buf = append(this.buf, level...)
	this.buf = append(this.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		this.buf = append(this.buf, '\n')
	}
	_, err := this.writeCloser.Write(this.buf)
	return err
}

func (this *Logger) println(level string, content ...interface{}) {
	if this.register.getLType() == FileLoggerType {
		this.rwLock.RLock()
		defer this.rwLock.RUnlock()
	}
	this.output(defaultCallDepth, level, fmt.Sprint(content...))
}

func (this *Logger) printf(level string, format string, content ...interface{}) {
	if this.register.getLType() == FileLoggerType {
		this.rwLock.RLock()
		defer this.rwLock.RUnlock()
	}
	this.output(defaultCallDepth, level, fmt.Sprintf(format, content...))
}

func (this *Logger) Info(content ...interface{}) {
	this.println(infoFlag, content...)
}

func (this *Logger) Infof(format string, content ...interface{}) {
	this.printf(infoFlag, format, content...)
}

func (this *Logger) Warn(content ...interface{}) {
	this.println(warnFlag, content...)
}

func (this *Logger) Warnf(format string, content ...interface{}) {
	this.printf(warnFlag, format, content...)
}

func (this *Logger) Error(content ...interface{}) {
	this.println(errorFlag, content...)
}

func (this *Logger) Errorf(format string, content ...interface{}) {
	this.printf(errorFlag, format, content...)
}

func (this *Logger) Crit(content ...interface{}) {
	this.println(criticalFlag, content...)
}

func (this *Logger) Critf(format string, content ...interface{}) {
	this.printf(criticalFlag, format, content...)
}

func Println(content ...interface{}) {
	defaultLogger.println(infoFlag, content...)
}

func Printf(format string, content ...interface{}) {
	defaultLogger.printf(infoFlag, format, content...)
}

func Info(content ...interface{}) {
	defaultLogger.println(infoFlag, content...)
}

func Infof(format string, content ...interface{}) {
	defaultLogger.printf(infoFlag, format, content...)
}

func Warn(content ...interface{}) {
	defaultLogger.println(warnFlag, content...)
}

func Warnf(format string, content ...interface{}) {
	defaultLogger.printf(warnFlag, format, content...)
}

func Error(content ...interface{}) {
	defaultLogger.println(errorFlag, content...)
}

func Errorf(format string, content ...interface{}) {
	defaultLogger.printf(errorFlag, format, content...)
}

func Crit(content ...interface{}) {
	defaultLogger.println(criticalFlag, content...)
}

func Critf(format string, content ...interface{}) {
	defaultLogger.printf(criticalFlag, format, content...)
}

func Fatal(content ...interface{}) {
	defaultLogger.println(fatalFlag, content...)
	os.Exit(1)
}

func Fatalf(format string, content ...interface{}) {
	defaultLogger.printf(fatalFlag, format, content...)
	os.Exit(1)
}

func Panic(content ...interface{}) {
	defaultLogger.println(panicFlag, content...)
	panic(fmt.Sprintln(content...))
}

func Panicf(format string, content ...interface{}) {
	defaultLogger.printf(panicFlag, format, content...)
	panic(fmt.Sprintf(format, content...))
}
