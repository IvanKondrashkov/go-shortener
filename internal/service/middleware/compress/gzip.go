package compress

import (
	"net/http"
	"strings"
)

// Header возвращает заголовки HTTP-ответа.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write сжимает данные и записывает их в ответ.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader устанавливает заголовок "Content-Encoding: gzip" и статус ответа.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < StatusCode {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// Read читает и распаковывает данные.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрывает gzip.Reader и исходный io.ReadCloser.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Gzip возвращает middleware для сжатия ответов и распаковки запросов в формате gzip.
func Gzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		ok := strings.Contains(acceptEncoding, "gzip")
		if ok {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		ok = strings.Contains(contentEncoding, "gzip")
		if ok {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Decompress is incorrect!"))
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}
