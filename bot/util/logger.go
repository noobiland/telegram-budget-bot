package util


import (
    "log/slog"
    "os"
)

var Logger *slog.Logger

func init() {
    Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelError, // Only log errors and higher levels
    }))
}