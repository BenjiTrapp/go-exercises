package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func MetaImageSearch(q string, timeout time.Duration) ([]Image, Stats, error) {
	var wg sync.WaitGroup
	imgs := []Image{}
	stats := Stats{}

	finished := make(chan bool)
	image := make(chan Image)
	errC := make(chan error)

	getImages := func(url string, parser func(body io.Reader) []Image) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				errC <- fmt.Errorf("error getting '%s': %s", url, err)
				return
			}
			for _, img := range parser(resp.Body) {
				image <- img
				<-time.After(10 * time.Millisecond)
			}
		}()
	}

	getImages(UnsplashURL(q), ParseUnsplashResult)
	getImages(PixabayURL(q), ParsePixabayResult)
	getImages(PexelsURL(q), ParsePexelsResult)

	go func() {
		wg.Wait()
		finished <- true
	}()

	for {
		select {
		case <-time.After(timeout):
			return imgs, stats, fmt.Errorf("Timeout überschritten -- breche Suche ab.")
		case newImage := <-image:
			imgs = append(imgs, newImage)
			stats.Inc(newImage.Source)
		case err := <-errC:
			return imgs, stats, err
		case <-finished:
			return imgs, stats, nil
		}
	}

	return nil, nil, fmt.Errorf("uhm ... something weird happened")
}
