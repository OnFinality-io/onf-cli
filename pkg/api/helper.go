package api

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

func sign(data, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	_, _ = mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func contentChecksum(data interface{}) string {
	b, _ := json.Marshal(data)
	h := md5.New()
	_, _ = io.WriteString(h, string(b))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func canonicalize(header http.Header, prefix string) string {
	var data []string
	for key, _ := range header {
		k := strings.ToLower(key)
		if strings.HasPrefix(k, prefix) {
			data = append(data, strings.Join([]string{k, header.Get(key)}, ":"))
		}
	}
	sort.Strings(data)
	str := strings.Join(data, "\n")
	return base64.StdEncoding.EncodeToString([]byte(str))
}
