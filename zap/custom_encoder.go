package zap

import (
	"encoding/json"
	"fmt"

	"github.com/gw123/glog/common"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// customConsoleEncoder is a custom encoder that places trace_id in a fixed position
type customConsoleEncoder struct {
	zapcore.Encoder
	cfg     zapcore.EncoderConfig
	traceID string
}

func newCustomConsoleEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &customConsoleEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg),
		cfg:     cfg,
		traceID: "",
	}
}

func (enc *customConsoleEncoder) Clone() zapcore.Encoder {
	return &customConsoleEncoder{
		Encoder: enc.Encoder.Clone(),
		cfg:     enc.cfg,
		traceID: enc.traceID,
	}
}

// AddString implements ObjectEncoder interface to intercept trace_id
func (enc *customConsoleEncoder) AddString(key, val string) {
	if key == common.KeyTraceID {
		enc.traceID = val
		return
	}
	enc.Encoder.AddString(key, val)
}

func (enc *customConsoleEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// Extract trace_id from fields
	for _, field := range fields {
		if field.Key == common.KeyTraceID && field.Type == zapcore.StringType {
			enc.traceID = field.String
		}
	}

	// Don't filter out any fields - keep them all for display

	// Determine trace ID to display
	traceID := "[]"
	if enc.traceID != "" {
		traceID = "[" + enc.traceID + "]"
	}

	// Build custom entry with trace_id in fixed position
	buf := buffer.NewPool().Get()

	// Time
	buf.AppendString("[")
	buf.AppendString(entry.Time.Format(DateTimeFormat))
	buf.AppendString("]")
	buf.AppendString(" ")

	// Level
	buf.AppendString("[")
	buf.AppendString(entry.Level.String())
	buf.AppendString("]")
	buf.AppendString(" ")

	// Logger name
	if entry.LoggerName != "" {
		names := []rune(entry.LoggerName)
		if len(names) > 0 && names[0] == '-' {
			// Check if it's just "-" or has a suffix like "-.something"
			if len(entry.LoggerName) > 2 && entry.LoggerName[1] == '.' {
				// Has a named suffix, extract it
				buf.AppendString("[")
				buf.AppendString(entry.LoggerName[2:]) // Skip "-."
				buf.AppendString("]")
			} else {
				// Just "-", show empty brackets
				buf.AppendString("[]")
			}
		} else {
			buf.AppendString("[")
			// Remove first part before first dot
			loggerName := entry.LoggerName
			for i, c := range loggerName {
				if c == '.' {
					loggerName = loggerName[i+1:]
					break
				}
			}
			buf.AppendString(loggerName)
			buf.AppendString("]")
		}
		buf.AppendString(" ")
	}

	// Caller
	if entry.Caller.Defined {
		buf.AppendString(entry.Caller.TrimmedPath())
		buf.AppendString(" ")
	}

	// Trace ID - fixed position
	buf.AppendString(traceID)
	buf.AppendString(" ")

	// Message
	buf.AppendString(" ")
	buf.AppendString(entry.Message)

	// Other fields
	if len(fields) > 0 {
		buf.AppendString(" ")
		buf.AppendByte('{')
		for i, field := range fields {
			if i > 0 {
				buf.AppendString(", ")
			}
			buf.AppendByte('"')
			buf.AppendString(field.Key)
			buf.AppendString(`": `)

			switch field.Type {
			case zapcore.StringType:
				buf.AppendByte('"')
				buf.AppendString(field.String)
				buf.AppendByte('"')
			case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
				fmt.Fprintf(buf, "[%d]", field.Integer)
			case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type:
				fmt.Fprintf(buf, "%d", field.Integer)
			case zapcore.Float64Type, zapcore.Float32Type:
				fmt.Fprintf(buf, "%v", field.Integer)
			case zapcore.BoolType:
				if field.Integer == 1 {
					buf.AppendString("true")
				} else {
					buf.AppendString("false")
				}
			default:
				// Try to marshal as JSON for complex types
				if data, err := json.Marshal(field.Interface); err == nil {
					buf.Write(data)
				} else {
					fmt.Fprintf(buf, `"%v"`, field.Interface)
				}
			}
		}
		buf.AppendByte('}')
	}

	buf.AppendString("\n")
	return buf, nil
}
