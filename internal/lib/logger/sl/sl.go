package sl

import "log/slog"

//learn deeper about slog to make PrettySLog

func Err(err error) slog.Attr {
	return slog.Attr{
		Key: "error",
		Value: slog.StringValue(err.Error()),
	}
	
}
