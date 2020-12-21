package glog

import "testing"

func TestError(t *testing.T) {
	type args struct {
		format string
		params []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "abc",
			args: args{
				format: "this is a log %s",
				params: []string{"test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.format)
			Infof(tt.args.format, tt.args.params[0])
			Warn(tt.args.format)
			Warnf(tt.args.format, tt.args.params[0])
			Error(tt.args.format)
			Errorf(tt.args.format, tt.args.params[0])
			Debug(tt.args.format)
			Debugf(tt.args.format, tt.args.params[0])

			Info(tt.args.format)
			Infof(tt.args.format, tt.args.params[0])
			Warn(tt.args.format)
			Warnf(tt.args.format, tt.args.params[0])
			Error(tt.args.format)
			Errorf(tt.args.format, tt.args.params[0])
			Debug(tt.args.format)
			Debugf(tt.args.format, tt.args.params[0])
		})
	}
}
