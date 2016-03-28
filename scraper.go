package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"

	"text/template"
)

// Image represents an image that gets returned by an image search.
//
//   URL:    url of the image
//   Source: name of the search engine that yielded this image
type Image struct {
	URL    string
	Source string
}

func (i Image) String() string {
	return fmt.Sprintf("%15s: %s", i.Source, i.URL)
}

const (
	pexelURLFormatString    = "https://www.pexels.com/de-de/suche/%s/"
	pixabayURLFormatString  = "https://pixabay.com/images/search/%s/"
	unsplashURLFormatString = "https://unsplash.com/search/photos/%s/"
)

// PexelsURL returns the bing url that searches for a specific image
//
//   q: search term
func PexelsURL(q string) string {
	return fmt.Sprintf(pexelURLFormatString, template.URLQueryEscaper(q))
}

// PixabayURL returns the pixabay url that searches for a specific image
//
//   q: search term
func PixabayURL(q string) string {
	return fmt.Sprintf(pixabayURLFormatString, template.URLQueryEscaper(q))
}

// UnsplashURL returns the unsplash url that searches for a specific image
//
//   q: search term
func UnsplashURL(q string) string {
	return fmt.Sprintf(unsplashURLFormatString, template.URLQueryEscaper(q))
}

func extractImages(r io.Reader, urlRegexp string, subRegexp *string, sub *string, source string) []Image {
	images := []Image{}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	body := string(b)

	if subRegexp != nil {
		body = regexp.MustCompile(*subRegexp).ReplaceAllString(body, *sub)
	}
	res := regexp.MustCompile(urlRegexp).FindAllStringSubmatch(body, -1)

imageAddLoop:
	for _, url := range res {
		for _, img := range images {
			if img.URL == url[1] {
				continue imageAddLoop
			}
		}
		images = append(images, Image{URL: url[1], Source: source})
	}
	return images
}

// ParsePexelsResult parses a pexel response and returns an array of the images the search returned
func ParsePexelsResult(r io.Reader) []Image {
	//return extractImageUrls(r, "https://images.pexels.com/photos", "pexel")
	return extractImages(r, ` src=['"](https://images.pexels.com/photos.*?)/?['"]`, nil, nil, "pexels")
}

// ParsePixabayResult parses a pixabay response and returns an array of the images the search returned
func ParsePixabayResult(r io.Reader) []Image {
	return extractImages(r, ` src=['"](https://cdn.pixabay.com/photo/.*?)/?['"]`, nil, nil, "pixabay")
}

// ParseUnsplashResult parses a unsplash response and returns an array of the images the search returned
func ParseUnsplashResult(r io.Reader) []Image {
	subre := `\\u002F`
	sub := "/"
	return extractImages(r, `"small":['"](https://images.unsplash.com/photo.*?)/?['"]`, &subre, &sub, "unsplash")
}
