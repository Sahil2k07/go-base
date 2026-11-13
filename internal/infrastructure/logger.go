package infrastructure

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/utils"
)

func ConfigureLogging(fileName string) {
	var output io.Writer = os.Stdout

	appConfig := configs.GetAppConfig()

	if appConfig.FileLogging {
		logPath := filepath.Join(
			utils.ResolveProjectRoot(),
			"logs",
			fileName+".log",
		)

		os.MkdirAll(filepath.Dir(logPath), 0755)

		file, err := os.OpenFile(
			logPath,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0666,
		)
		if err != nil {
			slog.Error("Error opening log file for writing logs", "error", err)
			panic("")
		}

		output = io.MultiWriter(os.Stdout, file)
	}

	if appConfig.JSONLogs {
		handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
				if attr.Key == slog.TimeKey {
					attr.Value = slog.TimeValue(attr.Value.Time().UTC())
				}
				return attr
			},
		})

		logger := slog.New(handler)
		slog.SetDefault(logger)

		log.SetOutput(
			slog.NewLogLogger(logger.Handler(), slog.LevelInfo).Writer(),
		)
		return
	}

	log.SetOutput(output)
}
