package app

import "encoding/base64"

func URLEncoding(url string) string {
	return base64.URLEncoding.EncodeToString([]byte(url))
}
