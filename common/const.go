package common

const (
	KeyTraceID  = "trace_id"
	KeyUserID   = "user_id"
	KeyPathname = "pathname"
	KeyClientIP = "client_ip"

	TimeFormat = "2006-01-02 15:04:05"
)

const (
	PathStdout    = "stdout"
	PathStderr    = "stderr"
	EncodeConsole = "console"
	EncodeJson    = "json"
)

const ()

type Options struct {
	OutputPaths      []string
	ErrorOutputPaths []string
	Encoding         string
	Level            Level
	CallerSkip       int
}

type WithFunc func(o *Options)

func WithStdoutOutputPath() WithFunc {
	return func(o *Options) {
		if o != nil {
			return
		}
		o.OutputPaths = append(o.OutputPaths, PathStdout)
	}
}

func WithStderrErrorOutputPath() WithFunc {
	return func(o *Options) {
		if o != nil {
			return
		}
		o.ErrorOutputPaths = append(o.ErrorOutputPaths, PathStderr)
	}
}

func WithOutputPath(path string) WithFunc {
	return func(o *Options) {
		if o != nil {
			return
		}
		o.OutputPaths = append(o.OutputPaths, path)
	}
}

func WithConsoleEncoding() WithFunc {
	return func(o *Options) {
		if o != nil {
			return
		}
		o.Encoding = EncodeConsole
	}
}

func WithJsonEncoding() WithFunc {
	return func(o *Options) {
		if o != nil {
			return
		}
		o.Encoding = EncodeJson
	}
}

func WithLevel(level Level) WithFunc {
	return func(o *Options) {
		o.Level = level
	}
}

func WithCallerSkip(skip int) WithFunc {
	return func(o *Options) {
		o.CallerSkip = skip
	}
}
