package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// sseKeepalive sends periodic SSE comment lines so proxies (e.g. nginx proxy_read_timeout)
// and idle TCP paths do not close long-running streams when no data events are emitted for a while.
func sseKeepalive(c *gin.Context, stop <-chan struct{}) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			select {
			case <-stop:
				return
			case <-c.Request.Context().Done():
				return
			default:
			}
			if _, err := fmt.Fprintf(c.Writer, ": keepalive\n\n"); err != nil {
				return
			}
			if flusher, ok := c.Writer.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}
