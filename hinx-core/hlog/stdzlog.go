// @Title stdzlog.go
// @Description Wraps zlogger log methods to provide global methods
// @Author Aceld - Thu Mar 11 10:32:29 CST 2019
package hlog

/*
	A global Log handle is provided by default for external use, which can be called directly through the API series.
	The global log object is StdHinxLog.
	Note: The methods in this file do not support customization and cannot replace the log recording mode.

	If you need a custom logger, please use the following methods:
	zlog.SetLogger(yourLogger)
	zlog.Ins().InfoF() and other methods.

   全局默认提供一个Log对外句柄，可以直接使用API系列调用
   全局日志对象 StdHinxLog
   注意：本文件方法不支持自定义，无法替换日志记录模式，如果需要自定义Logger:

   请使用如下方法:
   zlog.SetLogger(yourLogger)
   zlog.Ins().InfoF()等方法
*/

// StdHinxLog creates a global log
var StdHinxLog = NewHinxLog("", BitDefault)

// Flags gets the flags of StdHinxLog
func Flags() int {
	return StdHinxLog.Flags()
}

// ResetFlags sets the flags of StdHinxLog
func ResetFlags(flag int) {
	StdHinxLog.ResetFlags(flag)
}

// AddFlag adds a flag to StdHinxLog
func AddFlag(flag int) {
	StdHinxLog.AddFlag(flag)
}

// SetPrefix sets the log prefix of StdHinxLog
func SetPrefix(prefix string) {
	StdHinxLog.SetPrefix(prefix)
}

// SetLogFile sets the log file of StdHinxLog
func SetLogFile(fileDir string, fileName string) {
	StdHinxLog.SetLogFile(fileDir, fileName)
}

// SetMaxAge 最大保留天数
func SetMaxAge(ma int) {
	StdHinxLog.SetMaxAge(ma)
}

// SetMaxSize 单个日志最大容量 单位：字节
func SetMaxSize(ms int64) {
	StdHinxLog.SetMaxSize(ms)
}

// SetCons 同时输出控制台
func SetCons(b bool) {
	StdHinxLog.SetCons(b)
}

// SetLogLevel sets the log level of StdHinxLog
func SetLogLevel(logLevel int) {
	StdHinxLog.SetLogLevel(logLevel)
}

func Debugf(format string, v ...interface{}) {
	StdHinxLog.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	StdHinxLog.Debug(v...)
}

func Infof(format string, v ...interface{}) {
	StdHinxLog.Infof(format, v...)
}

func Info(v ...interface{}) {
	StdHinxLog.Info(v...)
}

func Warnf(format string, v ...interface{}) {
	StdHinxLog.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	StdHinxLog.Warn(v...)
}

func Errorf(format string, v ...interface{}) {
	StdHinxLog.Errorf(format, v...)
}

func Error(v ...interface{}) {
	StdHinxLog.Error(v...)
}

func Fatalf(format string, v ...interface{}) {
	StdHinxLog.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	StdHinxLog.Fatal(v...)
}

func Panicf(format string, v ...interface{}) {
	StdHinxLog.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	StdHinxLog.Panic(v...)
}

func Stack(v ...interface{}) {
	StdHinxLog.Stack(v...)
}

func init() {
	// Since the StdHinxLog object wraps all output methods with an extra layer, the call depth is one more than a normal logger object
	// The call depth of a regular zinxLogger object is 2, and the call depth of StdHinxLog is 3
	// (因为StdHinxLog对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的logger对象多一层调用
	// 一般的zinxLogger对象 calldDepth=2, StdHinxLog的calldDepth=3)
	StdHinxLog.calldDepth = 3
}
