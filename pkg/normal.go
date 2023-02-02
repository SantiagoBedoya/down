package down

import (
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

func (a *App) Normal() error {
	req, err := http.NewRequest("GET", a.URI, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(a.Destination, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		a.Destination,
	)

	_, err = io.Copy(io.MultiWriter(bar, f), resp.Body)
	return err
}
