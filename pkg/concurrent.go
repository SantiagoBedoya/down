package down

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func (a *App) Concurrent() error {
	fileSize, err := a.getHeaderDetail()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	sizePerWorker := fileSize / a.Concurrency

	bar := progressbar.DefaultBytes(
		int64(fileSize),
		a.Destination,
	)

	for i := 0; i < a.Concurrency; i++ {
		wg.Add(1)
		a.Chunks[i] = make([]byte, 0)
		start := i * sizePerWorker
		dataRange := fmt.Sprintf("bytes=%d-%d", start, start+sizePerWorker-1)

		if i == a.Concurrency-1 {
			dataRange = fmt.Sprintf("bytes=%d-", start)
		}
		go a.download(wg, i, dataRange, bar)
	}

	wg.Wait()
	if a.Err != nil {
		return a.Err
	}

	return a.combineChunks()
}

func (a *App) combineChunks() error {
	f, err := os.OpenFile(a.Destination, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	for i := 0; i < len(a.Chunks); i++ {
		buf.Write(a.Chunks[i])
	}
	_, err = buf.WriteTo(f)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) getHeaderDetail() (int, error) {
	req, _ := http.NewRequest("HEAD", a.URI, nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		log.Println("Error: StatusCode: ", resp.StatusCode)
		os.Exit(1)
	}

	header, ok := resp.Header["Content-Length"]
	if !ok {
		return 0, errors.New("file content length could not be determined")
	}
	fileSize, _ := strconv.Atoi(header[0])
	return fileSize, nil
}

func (a *App) download(wg *sync.WaitGroup, index int, dataRange string, bar *progressbar.ProgressBar) {
	defer wg.Done()

	req, err := http.NewRequest("GET", a.URI, nil)
	if err != nil {
		a.Lock()
		a.Err = err
		a.Unlock()
		return
	}
	req.Header.Add("Range", dataRange)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.Lock()
		a.Err = err
		a.Unlock()
		return
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)

	io.Copy(io.MultiWriter(bar, buf), resp.Body)

	a.Lock()
	a.Chunks[index] = append(a.Chunks[index], buf.Bytes()...)
	a.Unlock()
	io.Copy(bar, resp.Body)
}
