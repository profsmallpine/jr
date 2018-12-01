package minify

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	javascriptTag = `<script src="/%s"></script>`
	stylesheetTag = `<link type="text/css" rel="stylesheet" href="/%s">`
)

// AssetHelperFuncs is package var to pass to static HTML
var AssetHelperFuncs = template.FuncMap{
	"javascriptTag": JavascriptTag,
	"stylesheetTag": StylesheetTag,
}

// JavascriptTag returns html with script tags for all files from assets.go in
// development and a single bundle file outside of development.
func JavascriptTag() template.HTML {
	assetURL := "assets/js/application-*"
	paths, mtimes := resolveAssetUrls(assetURL)
	return generateRawHTML(paths, mtimes, javascriptTag)
}

// StylesheetTag returns html with style tags for all files from assets.go in
// development and a single bundle file outside of development.
func StylesheetTag() template.HTML {
	assetURL := "assets/css/application-*"
	paths, mtimes := resolveAssetUrls(assetURL)
	return generateRawHTML(paths, mtimes, stylesheetTag)
}

func resolveAssetUrls(assetURL string) (urls []string, mtimes []time.Time) {
	env := os.Getenv("ENVIRONMENT")
	if env == "test" {
		return []string{}, []time.Time{}
	}
	if env == "development" {
		return getUnbundledAssets(assetURL)
	}
	return getBundledAssets(assetURL)
}

func getUnbundledAssets(assetURL string) (urls []string, mtimes []time.Time) {
	if strings.Contains(assetURL, "js") {
		urls = javascripts
	} else {
		urls = stylesheets
	}

	for _, assetPath := range urls {
		info, err := os.Stat(assetPath)
		if err != nil {
			log.Fatalln(err)
		}
		mtimes = append(mtimes, info.ModTime())
	}
	return
}

func getBundledAssets(assetURL string) (urls []string, mtimes []time.Time) {
	fls, err := filepath.Glob(assetURL)
	if err != nil || fls == nil {
		return getUnbundledAssets(assetURL)
	}
	urls = []string{fls[0]}
	mtimes = nil

	return
}

func generateRawHTML(urls []string, mtimes []time.Time, tag string) template.HTML {
	htmls := []string{}

	for i, url := range urls {
		murl := url
		if mtimes != nil {
			murl += "?" + strconv.FormatInt(mtimes[i].Unix(), 10)
		}
		htmls = append(htmls, fmt.Sprintf(tag, murl))
	}
	return template.HTML(strings.Join(htmls, "\n"))
}
