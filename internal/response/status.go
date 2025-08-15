package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	OK                    StatusCode = 200
	BAD_REQUEST           StatusCode = 400
	INTERNAL_SERVER_ERROR StatusCode = 500
)

const BLANK_RESPONSE_REASON = " "

var statusCodeReasons = map[StatusCode]string{
	OK:                    " OK",
	BAD_REQUEST:           " Bad Request",
	INTERNAL_SERVER_ERROR: " Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write([]byte(statusLine(statusCode)))
	return err
}

func statusLine(statusCode StatusCode) string {
	return fmt.Sprintf("HTTP/1.1 %d%s\r\n", statusCode, statusReason(statusCode))
}

func statusReason(statusCode StatusCode) string {
	if response, ok := statusCodeReasons[statusCode]; ok {
		return response
	} else {
		return BLANK_RESPONSE_REASON
	}
}

func HtmlHead(StatusCode StatusCode) string {
	return fmt.Sprintf(`<head>
    <title>%d%s</title>
  </head>`, StatusCode, statusReason(StatusCode))
}

func htmlBody(h1 string, p string) string {
	return fmt.Sprintf(`<body>
    <h1>%s</h1>
	<p>%s</p>
  </body>`, h1, p)

}

func HtmlHandlerMessage(statusCode StatusCode, h1 string, p string) string {
	return fmt.Sprintf(`<html>
  %s
  %s
</html>`, HtmlHead(statusCode), htmlBody(h1, p))
}

func HtmlHandlerErrorMessage(statusCode StatusCode, p string) string {
	return HtmlHandlerMessage(statusCode, statusReason(statusCode), p)
}
