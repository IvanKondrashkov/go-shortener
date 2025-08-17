package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

// StatusCode определяет порог статуса, при котором сжатие ответа не применяется.
const (
	StatusCode = 300
)

// compressWriter реализует http.ResponseWriter для сжатия ответа в формате gzip.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// compressReader реализует io.ReadCloser для распаковки gzip-сжатых данных.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// newCompressWriter создает новый compressWriter.
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// newCompressReader создает новый compressReader.
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}
