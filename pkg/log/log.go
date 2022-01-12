// Copyright 2018-2022 (c) The Go Service Components authors. All rights reserved. Issued under the Apache 2.0 License.

package log // import "github.com/karlmutch/go-service/pkg/log"

// This file contains the implementation of a logger that adorns the logxi package with
// some common information not by default supplied by the generic code

import (
	"os"
	"sync"

	"github.com/go-stack/stack"
	logxi "github.com/karlmutch/logxi/v1"
)

var (
	hostName string
)

func init() {
	hostName, _ = os.Hostname()
}

// Logger encapsulates the logging device that is used to emit logs and
// as a receiver that has the logging methods
//
type Logger struct {
	log        logxi.Logger // The base implementation that is being encapsulated
	debugStack bool         // Should a debug stack be produced with the message
	hostName   string       // An optional host name to be used with the message
	sync.Mutex
}

// NewLogger can be used to instantiate a wrapper logger with a module label with
// output going stdout
//
func NewLogger(component string) (log *Logger) {
	logxi.DisableCallstack()

	return &Logger{
		log:        logxi.New(component),
		hostName:   hostName,
		debugStack: true,
	}
}

// NewErrLogger can be used to instantiate a wrapper logger with a module label with
// output going stderr
//
func NewErrLogger(component string) (log *Logger) {
	logxi.DisableCallstack()

	return &Logger{
		log:        logxi.NewLogger(logxi.NewConcurrentWriter(os.Stderr), component),
		hostName:   hostName,
		debugStack: true,
	}
}

// IncludeStack is used to enable a small function call stack to be included with messages
//
func (l *Logger) IncludeStack(included bool) (log *Logger) {
	l.Lock()
	defer l.Unlock()

	l.debugStack = included
	return l
}

// HostName is used to add an optional host name to messages, if empty then the host name will not be output
//
func (l *Logger) HostName(hostName string) (log *Logger) {
	l.Lock()
	defer l.Unlock()

	l.hostName = hostName
	return l
}

// Trace is a method for output of trace level messages
// with a varargs style list of parameters that is formatted
// as label and then the value in a single list
//
func (l *Logger) Trace(msg string, args ...interface{}) {
	if !l.IsTrace() {
		return
	}

	allArgs := append([]interface{}{}, args...)

	l.Lock()
	defer l.Unlock()

	if l.debugStack {
		allArgs = append(allArgs, "stack")
		allArgs = append(allArgs, stack.Trace()[1:].TrimRuntime())
	}

	if len(l.hostName) != 0 {
		allArgs = append(allArgs, "host")
		allArgs = append(allArgs, hostName)
	}

	l.log.Trace(msg, allArgs)
}

// Debug is a method for output of debugging level messages
// with a varargs style list of parameters that is formatted
// as label and then the value in a single list
//
func (l *Logger) Debug(msg string, args ...interface{}) {
	if !l.IsDebug() {
		return
	}

	allArgs := append([]interface{}{}, args...)

	l.Lock()
	defer l.Unlock()

	if l.debugStack {
		allArgs = append(allArgs, "stack")
		allArgs = append(allArgs, stack.Trace()[1:].TrimRuntime())
	}

	if len(l.hostName) != 0 {
		allArgs = append(allArgs, "host")
		allArgs = append(allArgs, hostName)
	}

	l.log.Debug(msg, allArgs)
}

// Info is a method for output of informational level messages
// with a varargs style list of parameters that is formatted
// as label and then the value in a single list
//
func (l *Logger) Info(msg string, args ...interface{}) {
	if !l.IsInfo() {
		return
	}

	allArgs := append([]interface{}{}, args...)

	l.Lock()
	defer l.Unlock()

	if l.debugStack {
		allArgs = append(allArgs, "stack")
		allArgs = append(allArgs, stack.Trace()[1:].TrimRuntime())
	}

	if len(l.hostName) != 0 {
		allArgs = append(allArgs, "host")
		allArgs = append(allArgs, hostName)
	}

	l.log.Info(msg, allArgs)
}

// Warn is a method for output of warning level messages
// with a varargs style list of parameters that is formatted
// as label and then the value in a single list
//
func (l *Logger) Warn(msg string, args ...interface{}) error {
	if !l.IsWarn() {
		return nil
	}

	allArgs := append([]interface{}{}, args...)

	l.Lock()
	defer l.Unlock()

	if l.debugStack {
		allArgs = append(allArgs, "stack")
		allArgs = append(allArgs, stack.Trace()[1:].TrimRuntime())
	}

	if len(l.hostName) != 0 {
		allArgs = append(allArgs, "host")
		allArgs = append(allArgs, hostName)
	}

	return l.log.Warn(msg, allArgs)
}

// Error is a method for output of error level messages
// with a varargs style list of parameters that is formatted
// as label and then the value in a single list
//
func (l *Logger) Error(msg string, args ...interface{}) error {

	allArgs := append([]interface{}{}, args...)

	l.Lock()
	defer l.Unlock()

	allArgs = append(allArgs, "stack")
	allArgs = append(allArgs, stack.Trace()[1:].TrimRuntime())

	if len(l.hostName) != 0 {
		allArgs = append(allArgs, "host")
		allArgs = append(allArgs, hostName)
	}

	return l.log.Error(msg, allArgs)
}

// Fatal is a method for output of fatal level messages
// with a varargs style list of parameters that is formatted
// as label and then the value in a single list
//
func (l *Logger) Fatal(msg string, args ...interface{}) {
	allArgs := append([]interface{}{}, args...)

	l.Lock()
	defer l.Unlock()

	allArgs = append(allArgs, "stack")
	allArgs = append(allArgs, stack.Trace()[1:].TrimRuntime())

	if len(l.hostName) != 0 {
		allArgs = append(allArgs, "host")
		allArgs = append(allArgs, hostName)
	}

	l.log.Fatal(msg, allArgs)
}

// Log is a method for output of parameterized level messages
// with a varargs style list of parameters that is formatted
// as label and then the value in a single list
//
func (l *Logger) Log(level int, msg string, args []interface{}) {
	allArgs := append([]interface{}{}, args...)

	l.Lock()
	defer l.Unlock()

	if level < logxi.LevelWarn {
		allArgs = append(allArgs, "stack")
		allArgs = append(allArgs, stack.Trace()[1:].TrimRuntime())
	}

	if len(l.hostName) != 0 {
		allArgs = append(allArgs, "host")
		allArgs = append(allArgs, hostName)
	}

	l.log.Log(level, msg, allArgs)
}

// SetLevel can be used to set the threshold for the level of messages
// that will be output by the logger
//
func (l *Logger) SetLevel(lvl int) {
	l.Lock()
	defer l.Unlock()
	l.log.SetLevel(lvl)
}

// IsTrace returns true in the event that the theshold logging level
// allows for trace messages to appear in the output
//
func (l *Logger) IsTrace() bool {
	l.Lock()
	defer l.Unlock()
	return l.log.IsTrace()
}

// IsDebug returns true in the event that the theshold logging level
// allows for debugging messages to appear in the output
//
func (l *Logger) IsDebug() bool {
	l.Lock()
	defer l.Unlock()
	return l.log.IsDebug()
}

// IsInfo returns true in the event that the theshold logging level
// allows for informational messages to appear in the output
//
func (l *Logger) IsInfo() bool {
	l.Lock()
	defer l.Unlock()
	return l.log.IsInfo()
}

// IsWarn returns true in the event that the theshold logging level
// allows for warning messages to appear in the output
//
func (l *Logger) IsWarn() bool {
	l.Lock()
	defer l.Unlock()
	return l.log.IsWarn()
}
