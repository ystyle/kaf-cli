package analytics

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	analytics "github.com/ystyle/google-analytics"
	"github.com/ystyle/kaf-cli/internal/utils"
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
				Params: map[string]any{
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

func getClientID() string {
	clientID := fmt.Sprintf("%d", rand.Uint32())
	config, err := os.UserConfigDir()
	if err != nil {
		return clientID
	}
	dir := filepath.Join(config, "kaf-cli")
	filepath := filepath.Join(dir, "config")
	if exist, _ := utils.IsExists(filepath); exist {
		bs, err := os.ReadFile(filepath)
		if err != nil {
			return clientID
		}
		clientID = string(bs)
	} else {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return clientID
		}
		_ = os.WriteFile(filepath, []byte(clientID), 0700)
	}
	return clientID
}
