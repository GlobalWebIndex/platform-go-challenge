// Package logs provides a global logger for zerolog.
// Enhancement of github.com/rs/zerolog/log of Olivier Poitrey
package logs

import (
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"x-gwi/app/x/env"
)

//nolint:gochecknoinits
func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	zerolog.TimestampFieldName = "µs"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack //nolint:reassign
	zerolog.CallerMarshalFunc = formatCaller
	// zerolog.SetGlobalLevel(zerolog.Disabled)
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	cfg := &Logs{
		App: env.Env("APP_NAME", "app"),
		Enable: EnableOutput{
			SysLog:  false,
			StdErr:  false,
			StdOut:  false,
			Console: true,
			None:    false,
		},
	}
	cfg.SetOnce()
}

//nolint:gochecknoglobals
var (
	Logger   zerolog.Logger
	LogC2    zerolog.Logger // LogC2 global logger with caller in place - CallerWithSkipFrameCount(2)
	LogC3    zerolog.Logger // LogC3 global logger with caller one before - CallerWithSkipFrameCount(3)
	setOK    bool
	basepath string
)

// Logs configuration & credentials of logs // `yaml:"logs" json:"logs"`.
type Logs struct {
	App    string `json:"app" yaml:"app"`
	Enable EnableOutput
}

type EnableOutput struct {
	SysLog  bool `json:"syslog"  yaml:"syslog"`
	StdErr  bool `json:"stderr"  yaml:"stderr"`
	StdOut  bool `json:"stdout"  yaml:"stdout"`
	Console bool `json:"console" yaml:"console"`
	None    bool `json:"none"    yaml:"none"`
}

// explicit confirmation to not log on any writer
// (in pair with None).
func (c *Logs) anyEnabled() bool {
	return c.Enable.SysLog ||
		c.Enable.StdErr ||
		c.Enable.StdOut ||
		c.Enable.Console
}

func (c *Logs) setSyslog() io.Writer {
	if !c.Enable.SysLog {
		return zerolog.Nop()
	}

	//nolint:nosnakecase
	syslogWriter, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL6, c.App)
	if err != nil {
		Error().Err(err).Msg("Unable to setup syslog")

		return zerolog.Nop()
	}

	return zerolog.SyslogCEEWriter(syslogWriter)
}

func (c *Logs) setStdErr() io.Writer {
	if !c.Enable.StdErr {
		if c.anyEnabled() ||
			(!c.anyEnabled() && c.Enable.None) {
			return zerolog.Nop()
		}
	}

	return os.Stderr
}

func (c *Logs) setStdOut() io.Writer {
	if !c.Enable.StdOut {
		return zerolog.Nop()
	}

	return os.Stdout
}

func (c *Logs) setConsole() io.Writer {
	if !c.Enable.Console {
		return zerolog.Nop()
	}

	//nolint:exhaustruct
	return zerolog.ConsoleWriter{
		Out: os.Stdout,
		// StampMicro = "Jan _2 15:04:05.000000"
		TimeFormat: time.StampMicro,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			"app",
			"pid",
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
		PartsExclude: []string{
			"host",
		},
	}
}

// SetOnce() passes cfg at runtime.
//
//nolint:gomnd
func (c *Logs) SetOnce() {
	if setOK || c.App == "" {
		return
	}

	var mu sync.Mutex

	mu.Lock()
	defer mu.Unlock()

	basepath, _ = os.Getwd()
	Logger = zerolog.New(zerolog.MultiLevelWriter(
		c.setSyslog(),
		c.setStdErr(),
		c.setStdOut(),
		c.setConsole(),
	)).
		With().
		Timestamp().
		Str("app", c.App).
		Int("pid", os.Getpid()).
		Logger()

	LogC2 = Logger.With().
		CallerWithSkipFrameCount(2).
		// Int("pkg-logs", 1).
		Logger()

	LogC3 = Logger.With().
		CallerWithSkipFrameCount(3).
		// Int("pkg-logs", 1).
		Logger()

	log.SetFlags(0)
	log.SetOutput(Logger.With().
		CallerWithSkipFrameCount(5). // 4
		// Int("pkg-log", 1).
		Str("deprecated-log", "log").
		Logger())

	setOK = true
}

//nolint:nestif
func formatCaller(pc uintptr, file string, line int) string {
	if len(file) > 0 && len(basepath) > 0 {
		if rel, err := filepath.Rel(basepath, file); err == nil {
			file = filepath.Clean(rel)
			if strings.Contains(file, "src/") {
				file = file[strings.Index(file, "src/")+len("src/"):]
			}

			if strings.Contains(file, "../") {
				file = file[strings.Index(file, "../")+len("../"):]
			}
		}
	}
	// https://github.com/rs/zerolog/issues/456   .Str("func", runtime.FuncForPC(pc).Name())
	return file + ":" + strconv.Itoa(line) + " " + filepath.Base(runtime.FuncForPC(pc).Name()) + "()"
}

//nolint:gomnd
func LogTrace() (string, string) {
	pc := make([]uintptr, 10) // min 1
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	return "trace", fmt.Sprintf("%s:%d::%s()", filepath.Base(file), line, filepath.Base(f.Name()))
}

//nolint:gomnd,funlen,nestif
func LogTrace2() (string, string) {
	pc := make([]uintptr, 10) // min 1
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	bi, ok := debug.ReadBuildInfo()
	// bi, ok := debug.ReadBuildInfo()
	// f.
	basepath, _ := os.Getwd()
	if len(file) > 0 && len(basepath) > 0 {
		if rel, err := filepath.Rel(basepath, file); err == nil {
			file = filepath.Clean(rel)
			if strings.Contains(file, "src/") {
				file = file[strings.Index(file, "src/")+len("src/"):]
			}

			if strings.Contains(file, "../") {
				file = file[strings.Index(file, "../")+len("../"):]
			}
		}
	}

	fn := f.Name()
	bpm := bi.Path

	//nolint:godox
	// 	TODO: on test  bi. is empty
	// module
	if bpm == "" {
		bpm = "x-gwi"
	}

	if len(fn) > 0 && len(bpm) > 0 {
		if strings.Contains(fn, bpm) {
			fn = fn[strings.Index(fn, bpm+"/")+len(bpm+"/"):]
		}
	}
	//
	// fmt.Println(logs.LogTrace())
	// fmt.Println(logs.LogTrace2())
	// os.Exit(0)
	//
	// 	trace serv-grpc.go:24::serv.runServerGRCP()
	// trace serv-grpc.go:25::serv.runServerGRCP()
	// :/:github.com/ac-i/user-service/implem
	// :f:github.com/ac-i/user-service/implem/serv.runServerGRCP:
	// [true]
	// :bp:github.com/ac-i/user-service
	// :bmp:github.com/ac-i/user-service
	// [implem/serv.runServerGRCP]
	//
	// trace pki.go:28::pki.NewPKI()
	// trace pki.go:29::pki.NewPKI():/:xm-com/app:f:xm-com/app/pki.NewPKI:
	// [true]
	// :bp:xm-com
	// :bmp:xm-com
	// :bpm:xm-com
	// [app/pki.NewPKI]
	//
	// trace pki_test.go:13::pki_test.ExampleNewPKI()
	// trace pki_test.go:14::pki_test.ExampleNewPKI():/:xm-com/app:f:xm-com/app/pki_test.ExampleNewPKI:
	// [true]
	// :bp:
	// :bmp:
	// :bpm:xm-com
	// [app/pki_test.ExampleNewPKI]
	//
	return "trace", fmt.Sprintf("%s:%d::%s():/:%s:f:%s:\n[%v]\n:bp:%s\n:bmp:%s\n:bpm:%s\n[%s]", filepath.Base(file), line,
		filepath.Base(f.Name()),
		filepath.Dir(f.Name()),
		f.Name(),
		ok,
		bi.Path,
		bi.Main.Path,
		bpm,
		fn,
	)
}

func XID() xid.ID {
	// github.com/rs/xid guid := xid.New() xid: 12 bytes, 20 chars, configuration free, sortable
	// guid := xid.New() guid.Machine() guid.Pid() guid.Time() guid.Counter() _ = xid.NilID()
	return xid.New()
}
