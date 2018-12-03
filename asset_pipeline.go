package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

// minifyAssets compiles (minified + concatenated) assets into single hashed files
func minifyAssets() bool {
	if success := minifyConcat(&stylesheets, "text/css", "assets/css/application.css"); !success {
		return false
	}

	if success := minifyConcat(&javascripts, "text/javascript", "assets/js/application.js"); !success {
		return false
	}

	return true
}

func minifyConcat(assets *[]string, mimetype, location string) bool {
	ars, success := loadAssets(assets, mimetype)
	if !success {
		return false
	}
	r := io.MultiReader(ars...)

	fw, success := openSaveLocation(location)
	if !success {
		for _, ar := range ars {
			ar.(io.ReadCloser).Close()
		}
		return false
	}
	w := bufio.NewWriter(fw)

	// Remove old file if it exists
	removeOldFile(location)

	// Minify the MultiReader files, close ars, flush buffer, and close fw
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.Minify(mimetype, w, r)
	for _, ar := range ars {
		ar.(io.ReadCloser).Close()
	}
	w.Flush()

	fw.Close()

	// Rename file with hash (e.g. application-<hash>.css)
	renameNewFile(location)

	return true
}

func loadAssets(assets *[]string, mimetype string) ([]io.Reader, bool) {
	// Open all assets and put into slice
	ars := make([]io.Reader, len(*assets))
	for i, src := range *assets {
		ar, err := os.Open(src)
		if err != nil {
			fmt.Println("error in os.Open: ", err)
			for _, ar := range ars {
				ar.(io.ReadCloser).Close()
			}
			return nil, false
		}
		ars[i] = ar

		if i > 0 && mimetype == "text/javascript" {
			// prepend newline when concatenating JS files
			ars[i] = newPrependReader(ar, []byte("\n"))
		} else {
			ars[i] = ar
		}
	}
	return ars, true
}

func removeOldFile(location string) {
	ns := strings.Split(location, ".")
	fls, err := filepath.Glob(ns[0] + "-*")
	if err != nil || fls == nil {
		fmt.Println("error finding old file: ", err)
		return
	}

	if err := os.Remove(fls[0]); err != nil {
		fmt.Println("error removing old file: ", err)
	}
}

func renameNewFile(location string) {
	bytes, err := ioutil.ReadFile(location)
	if err != nil {
		fmt.Println("error reading new file: ", err)
		return
	}

	hash := bytesHash(bytes)
	ns := strings.Split(location, ".")
	if err := os.Rename(location, ns[0]+"-"+hash+"."+ns[1]); err != nil {
		fmt.Println("error renaming new file: ", err)
	}
}

func bytesHash(bytes []byte) string {
	sum := sha1.Sum(bytes)
	return hex.EncodeToString([]byte(sum[:]))
}

func openSaveLocation(location string) (*os.File, bool) {
	// Open file to save minified assets to
	var fw *os.File
	fw, err := os.OpenFile(location, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("error in os.OpenFile: ", err)
		return nil, false
	}
	return fw, true
}

type prependReader struct {
	io.ReadCloser
	prepend []byte
}

func newPrependReader(r io.ReadCloser, prepend []byte) *prependReader {
	return &prependReader{r, prepend}
}

func (r *prependReader) Read(p []byte) (int, error) {
	if r.prepend != nil {
		n := copy(p, r.prepend)
		if n != len(r.prepend) {
			return n, io.ErrShortBuffer
		}
		r.prepend = nil
		return n, nil
	}
	return r.ReadCloser.Read(p)
}

var javascripts = []string{
	"assets/js/jquery.js",
	"assets/js/jquery.easing.js",
	"assets/js/jquery-ui.min.js",
	"assets/js/bootstrap.min.js",
	"assets/js/jquery.flexslider.js",
	"assets/js/background-check.min.js",
	"assets/js/jquery.fitvids.js",
	"assets/js/jquery.viewportchecker.js",
	"assets/js/jquery.stellar.min.js",
	"assets/js/wow.min.js",
	"assets/js/jquery.colorbox-min.js",
	"assets/js/owl.carousel.min.js",
	"assets/js/isotope.pkgd.min.js",
	"assets/js/masonry.pkgd.min.js",
	"assets/js/imagesloaded.pkgd.min.js",
	"assets/js/jPushMenu.js",
	"assets/js/jquery.fs.tipper.min.js",
	"assets/js/mediaelement-and-player.min.js",
	"assets/js/jquery.validate.min.js",
	"assets/js/theme.js",
	"assets/js/navigation.js",
	"assets/js/jquery.mb.YTPlayer.min.js",
	"assets/js/jquery.singlePageNav.js",
	"assets/js/contact-form.js",
	"assets/js/map.js",
	"assets/js/TweenLite.min.js",
	"assets/js/EasePack.min.js",
	"assets/js/pollyfill.js",
}

var stylesheets = []string{
	"assets/css/style.css",
	"assets/css/navigation.css",
	"assets/css/flexslider.css",
	"assets/css/owl.carousel.css",
	"assets/css/mediaelementplayer.css",
	"assets/css/colorbox.css",
	"assets/css/jquery.fs.tipper.css",
	"assets/css/bootstrap.css",
	"assets/css/font-awesome.css",
	"assets/css/ionicons.css",
	"assets/css/jPushMenu.css",
	"assets/css/animate.css",
	"assets/css/jquery-ui.css",
	"assets/css/emoji.css",
}
