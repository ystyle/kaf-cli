package kafcli

import (
	"fmt"
	"github.com/ystyle/google-analytics"
	"math/rand"
	"runtime"
	"time"
)

func Analytics(version, secret, measurement, format string) {
	if secret == "" || measurement == "" {
		return
	}
	t := time.Now().Unix()
	analytics.SetKeys(secret, measurement) // // required
	payload := analytics.Payload{
		ClientID: fmt.Sprintf("%d.%d", rand.Int31(), t), // required
		UserID:   getClientID(),
		Events: []analytics.Event{
			{
				Name: "kaf_cli", // required
				Params: map[string]interface{}{
					"os":      runtime.GOOS,
					"arch":    runtime.GOARCH,
					"format":  format,
					"version": version,
				},
			},
		},
	}
	analytics.Send(payload)
}
