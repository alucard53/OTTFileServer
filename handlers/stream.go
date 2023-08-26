package handlers

import (
	"FileServer/files"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Stream struct {
	l *log.Logger
}

func NewStreamHandler(l *log.Logger) Stream {
	return Stream{l}
}

func (s Stream) handleErr(err error, w http.ResponseWriter) {
	s.l.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func getRoute(id string) (string, error) {
	route, err := files.FileList.Search(id)

	if err != nil {
		return "", err
	}

	return route, nil
}

func (s Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	route, err := getRoute(r.URL.Query().Get("q"))

	if err != nil {
		s.handleErr(err, w)
		return
	}

	file, err := os.Open(route)

	defer file.Close()

	stat, err := file.Stat()

	if err != nil {
		s.handleErr(err, w)
		return
	}

	rangeHeader := r.Header.Get("Range")

	if rangeHeader == "" {
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
		http.ServeContent(w, r, "", stat.ModTime(), file)
		return
	}

	rangeParts := strings.Split(rangeHeader, "=")
	if len(rangeParts) != 2 || rangeParts[0] != "bytes" {
		s.handleErr(files.RangeError, w)
		return
	}

	rangeValues := strings.Split(rangeParts[1], "-")
	if len(rangeValues) != 2 || rangeParts[0] != "bytes" {
		s.handleErr(files.RangeError, w)
		return
	}

	start, err := strconv.ParseInt(rangeValues[0], 10, 64)
	if err != nil {
		s.handleErr(files.RangeError, w)
		return
	}

	end := stat.Size() - 1

	if rangeValues[1] != "" {
		end, err = strconv.ParseInt(rangeValues[1], 10, 64)

		if err != nil {
			s.handleErr(files.RangeError, w)
			return
		}
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	w.Header().Set("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(stat.Size(), 10))

	_, err = file.Seek(start, 0)

	if err != nil {
		s.l.Println(err)
	}

	http.ServeContent(w, r, "", stat.ModTime(), file)
	w.WriteHeader(http.StatusPartialContent)
}
