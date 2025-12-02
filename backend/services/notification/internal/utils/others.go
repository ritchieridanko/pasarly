package utils

import (
	"encoding/base64"
	"net/url"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func URLWithToken(baseURL, path, token string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	u.Path = path

	q := u.Query()
	q.Set("token", token)

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func MIMEBase64(value string) string {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(value)) + "?="
}

func TraceErr(s trace.Span, err error, message string) {
	s.RecordError(err)
	s.SetStatus(codes.Error, message)
}
