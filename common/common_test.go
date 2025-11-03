package common

import (
	"testing"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{DPanicLevel, "dpanic"},
		{PanicLevel, "panic"},
		{FatalLevel, "fatal"},
		{Level(99), "Level(99)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.level.String()
			if got != tt.expected {
				t.Errorf("Level.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLevel_CapitalString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{DPanicLevel, "DPANIC"},
		{PanicLevel, "PANIC"},
		{FatalLevel, "FATAL"},
		{Level(99), "LEVEL(99)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.level.CapitalString()
			if got != tt.expected {
				t.Errorf("Level.CapitalString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLevel_MarshalText(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got, err := tt.level.MarshalText()
			if err != nil {
				t.Errorf("Level.MarshalText() error = %v", err)
				return
			}
			if string(got) != tt.expected {
				t.Errorf("Level.MarshalText() = %v, want %v", string(got), tt.expected)
			}
		})
	}
}

func TestLevel_UnmarshalText(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
		wantErr  bool
	}{
		{"debug", DebugLevel, false},
		{"DEBUG", DebugLevel, false},
		{"info", InfoLevel, false},
		{"INFO", InfoLevel, false},
		{"", InfoLevel, false},
		{"warn", WarnLevel, false},
		{"WARN", WarnLevel, false},
		{"error", ErrorLevel, false},
		{"ERROR", ErrorLevel, false},
		{"dpanic", DPanicLevel, false},
		{"DPANIC", DPanicLevel, false},
		{"panic", PanicLevel, false},
		{"PANIC", PanicLevel, false},
		{"fatal", FatalLevel, false},
		{"FATAL", FatalLevel, false},
		{"invalid", InfoLevel, true},
		{"unknown", InfoLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var l Level
			err := l.UnmarshalText([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Level.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && l != tt.expected {
				t.Errorf("Level.UnmarshalText() = %v, want %v", l, tt.expected)
			}
		})
	}
}

func TestLevel_UnmarshalText_Nil(t *testing.T) {
	var l *Level
	err := l.UnmarshalText([]byte("info"))
	if err == nil {
		t.Error("UnmarshalText on nil Level should return error")
	}
}

func TestLevel_Set(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
		wantErr  bool
	}{
		{"debug", DebugLevel, false},
		{"info", InfoLevel, false},
		{"warn", WarnLevel, false},
		{"error", ErrorLevel, false},
		{"invalid", InfoLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var l Level
			err := l.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Level.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && l != tt.expected {
				t.Errorf("Level.Set() = %v, want %v", l, tt.expected)
			}
		})
	}
}

func TestLevel_Get(t *testing.T) {
	l := InfoLevel
	got := l.Get()
	if got != InfoLevel {
		t.Errorf("Level.Get() = %v, want %v", got, InfoLevel)
	}
}

func TestLevel_Enabled(t *testing.T) {
	tests := []struct {
		baseLevel  Level
		checkLevel Level
		expected   bool
	}{
		{InfoLevel, DebugLevel, false},
		{InfoLevel, InfoLevel, true},
		{InfoLevel, WarnLevel, true},
		{InfoLevel, ErrorLevel, true},
		{DebugLevel, DebugLevel, true},
		{DebugLevel, InfoLevel, true},
		{ErrorLevel, WarnLevel, false},
		{ErrorLevel, InfoLevel, false},
		{ErrorLevel, ErrorLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.baseLevel.String()+"_"+tt.checkLevel.String(), func(t *testing.T) {
			got := tt.baseLevel.Enabled(tt.checkLevel)
			if got != tt.expected {
				t.Errorf("Level.Enabled(%v, %v) = %v, want %v",
					tt.baseLevel, tt.checkLevel, got, tt.expected)
			}
		})
	}
}

func TestOptions_WithFunctions(t *testing.T) {
	opts := Options{}

	WithStdoutOutputPath()(&opts)
	if len(opts.OutputPaths) != 1 || opts.OutputPaths[0] != PathStdout {
		t.Error("WithStdoutOutputPath() failed")
	}

	WithStderrErrorOutputPath()(&opts)
	if len(opts.ErrorOutputPaths) != 1 || opts.ErrorOutputPaths[0] != PathStderr {
		t.Error("WithStderrErrorOutputPath() failed")
	}

	WithOutputPath("test.log")(&opts)
	if len(opts.OutputPaths) != 2 || opts.OutputPaths[1] != "test.log" {
		t.Error("WithOutputPath() failed")
	}

	WithConsoleEncoding()(&opts)
	if opts.Encoding != EncodeConsole {
		t.Error("WithConsoleEncoding() failed")
	}

	WithJsonEncoding()(&opts)
	if opts.Encoding != EncodeJson {
		t.Error("WithJsonEncoding() failed")
	}

	WithLevel(DebugLevel)(&opts)
	if opts.Level != DebugLevel {
		t.Error("WithLevel() failed")
	}

	WithCallerSkip(2)(&opts)
	if opts.CallerSkip != 2 {
		t.Error("WithCallerSkip() failed")
	}
}

func TestOptions_WithFunctions_NilOptions(t *testing.T) {
	var opts *Options

	WithStdoutOutputPath()(opts)
	WithStderrErrorOutputPath()(opts)
	WithOutputPath("test.log")(opts)
	WithConsoleEncoding()(opts)
	WithJsonEncoding()(opts)
	WithLevel(DebugLevel)(opts)
	WithCallerSkip(1)(opts)
}

func TestConstants(t *testing.T) {
	if KeyTraceID != "trace_id" {
		t.Errorf("KeyTraceID = %v, want trace_id", KeyTraceID)
	}
	if KeyUserID != "user_id" {
		t.Errorf("KeyUserID = %v, want user_id", KeyUserID)
	}
	if KeyPathname != "pathname" {
		t.Errorf("KeyPathname = %v, want pathname", KeyPathname)
	}
	if KeyClientIP != "client_ip" {
		t.Errorf("KeyClientIP = %v, want client_ip", KeyClientIP)
	}

	if PathStdout != "stdout" {
		t.Errorf("PathStdout = %v, want stdout", PathStdout)
	}
	if PathStderr != "stderr" {
		t.Errorf("PathStderr = %v, want stderr", PathStderr)
	}

	if EncodeConsole != "console" {
		t.Errorf("EncodeConsole = %v, want console", EncodeConsole)
	}
	if EncodeJson != "json" {
		t.Errorf("EncodeJson = %v, want json", EncodeJson)
	}

	if TimeFormat != "2006-01-02 15:04:05" {
		t.Errorf("TimeFormat = %v, want 2006-01-02 15:04:05", TimeFormat)
	}
}

func BenchmarkLevel_String(b *testing.B) {
	l := InfoLevel
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.String()
	}
}

func BenchmarkLevel_Enabled(b *testing.B) {
	baseLevel := InfoLevel
	checkLevel := WarnLevel
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = baseLevel.Enabled(checkLevel)
	}
}

func BenchmarkLevel_UnmarshalText(b *testing.B) {
	text := []byte("info")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var l Level
		l.UnmarshalText(text)
	}
}
