package shared

import (
	"io"
	"net/http"
)

type MonitoredWriter struct {
	writer        http.ResponseWriter
	bytesWritten  int64
	onBytesWritten func(int64)
}

func NewMonitoredWriter(w http.ResponseWriter, onBytesWritten func(int64)) *MonitoredWriter {
	return &MonitoredWriter{
		writer:         w,
		onBytesWritten: onBytesWritten,
	}
}

func (mw *MonitoredWriter) Write(data []byte) (int, error) {
	n, err := mw.writer.Write(data)
	if n > 0 {
		mw.bytesWritten += int64(n)
		if mw.onBytesWritten != nil {
			mw.onBytesWritten(int64(n))
		}
	}
	return n, err
}

func (mw *MonitoredWriter) Header() http.Header {
	return mw.writer.Header()
}

func (mw *MonitoredWriter) WriteHeader(statusCode int) {
	mw.writer.WriteHeader(statusCode)
}

func (mw *MonitoredWriter) BytesWritten() int64 {
	return mw.bytesWritten
}

type MonitoredReader struct {
	reader        io.Reader
	bytesRead     int64
	onBytesRead   func(int64)
}

func NewMonitoredReader(r io.Reader, onBytesRead func(int64)) *MonitoredReader {
	return &MonitoredReader{
		reader:      r,
		onBytesRead: onBytesRead,
	}
}

func (mr *MonitoredReader) Read(data []byte) (int, error) {
	n, err := mr.reader.Read(data)
	if n > 0 {
		mr.bytesRead += int64(n)
		if mr.onBytesRead != nil {
			mr.onBytesRead(int64(n))
		}
	}
	return n, err
}

