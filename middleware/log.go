package middleware

import (
	"fmt"
	"hotels_api/api"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

// showLog middleware handler shows network data log info
func ShowLog(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := time.Now()

		sr := api.NewStatusHTTP(w)
		next.ServeHTTP(sr, req)

		statusCode := sr.StatusCode

		tw := tabwriter.NewWriter(os.Stdout, 28, 4, 1, ' ', tabwriter.Debug)
		fmt.Fprintf(tw, "%v\t [%d: %s]\t %v\t %s\t %s\t %s\n", t.Format("02.01.2006 15:04:05 MST"), statusCode, http.StatusText(statusCode), time.Since(t), req.RemoteAddr, req.Method, req.URL.String())
		tw.Flush()
	})
}
