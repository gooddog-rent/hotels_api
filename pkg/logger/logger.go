package logger

import (
	"fmt"
	"hotels_api/pkg/headers/status"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

// Log middleware handler shows network data log info
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := time.Now()

		sr := status.NewStatusHTTP(w)
		next.ServeHTTP(sr, req)

		statusCode := sr.StatusCode

		tw := tabwriter.NewWriter(os.Stdout, 28, 4, 1, ' ', tabwriter.Debug)

		layout := "02.01.2006 15:04:05"

		fmt.Fprintf(tw, "%v\t [%d: %s]\t %v\t %s\t %s\t %s\n", t.Format(layout), statusCode, http.StatusText(statusCode), time.Since(t), req.RemoteAddr, req.Method, req.URL.String())
		tw.Flush()
	})
}
