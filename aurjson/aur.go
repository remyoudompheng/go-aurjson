package aur

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"time"
)

const AURLocation = "http://aur.archlinux.org/rpc.php"

type AURResponseType struct {
	Type string `json:"type"`
}

type ErrorResponse struct {
	AURResponseType
	Msg string `json:"results"`
}

type InfoResponse struct {
	Type    string   `json:"type"`
	Results *PkgInfo `json:"results"`
}

type SearchResponse struct {
	Type    string    `json:"type"`
	Results []PkgInfo `json:"results"`
}

// Sort by name interface
func (s SearchResponse) Len() int           { return len(s.Results) }
func (s SearchResponse) Less(i, j int) bool { return s.Results[i].Name < s.Results[j].Name }
func (s SearchResponse) Swap(i, j int)      { s.Results[i], s.Results[j] = s.Results[j], s.Results[i] }

var _ sort.Interface = SearchResponse{}

type PkgInfo struct {
	Name           string
	Version        string
	Description    string
	Maintainer     *string
	ID             int64 `json:",string"`
	CategoryID     int64 `json:",string"`
	License        string
	NumVotes       int `json:",string"`
	OutOfDate      int `json:",string"`
	URL, URLPath   string
	FirstSubmitted *JSONTime `json:",string"`
	LastModified   *JSONTime `json:",string"`
}

func genericQuery(querytype, arg string, target interface{}) error {
	form := make(url.Values)
	form.Add("type", querytype)
	form.Add("arg", arg)
	url := AURLocation + "?" + form.Encode()
	resp, er := http.Get(url)
	if er != nil {
		return er
	}

	buf := bytes.NewBuffer(nil)
	_, er = io.CopyN(buf, resp.Body, 1e6)
	if er != nil && er != io.EOF {
		return er
	}
	jsonthings := buf.Bytes()

	var dummy AURResponseType
	er = json.Unmarshal(jsonthings, &dummy)

	switch {
	case dummy.Type == "error":
		var aur_err ErrorResponse
		er = json.Unmarshal(jsonthings, &aur_err)
		return fmt.Errorf("Error from AUR server: %s", aur_err.Msg)
	case er != nil:
		return fmt.Errorf("error in JSON parser: %s", er)
	default:
		er = json.Unmarshal(jsonthings, target)
		return er
	}
	panic("impossible")
}

func GetInfo(pkg string) (*PkgInfo, error) {
	var info InfoResponse
	er := genericQuery("info", pkg, &info)
	return info.Results, er
}

func DoSearch(pattern string) ([]PkgInfo, error) {
	var info SearchResponse
	er := genericQuery("search", pattern, &info)
	sort.Sort(info)
	return info.Results, er
}

type JSONTime time.Time

func (t *JSONTime) UnmarshalJSON(j []byte) error {
	var x int64
	if er := json.Unmarshal(j, &x); er != nil {
		return er
	}
	*t = JSONTime(time.Unix(x, 0))
	return nil
}

func (t *JSONTime) String() string {
	t2 := time.Time(*t)
	return t2.Format("2006-01-02 15:04:05 MST")
}
