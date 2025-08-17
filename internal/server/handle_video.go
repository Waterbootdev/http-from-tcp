package server

import (
	"os"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

func handleVideo(w *response.Writer) {

	video, err := os.ReadFile("assets/vim.mp4")

	if err != nil {
		internalServerError(err).write(w)
		return
	}

	if w.WriteStatusLine(response.OK) {
		return
	}
	if w.WriteHeaders(headers.GetContentTypeHeaders(len(video), headers.MP4)) {
		return
	}

	if w.WriteBody(video) {
		return
	}
}
