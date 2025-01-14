package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
)

type prettifyWriter struct {
	writer *os.File
}

func (p *prettifyWriter) Write(data []byte) (int, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return p.writer.Write(data)
	}

	prettyJSON.WriteByte('\n')
	return p.writer.Write(prettyJSON.Bytes())
}

func GetJsonLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(&prettifyWriter{os.Stdout}, nil))
}
