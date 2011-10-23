package aur

import (
	"os"
	"url"
	"http"
	"json"
)

const AURLocation = "http://aur.archlinux.org/rpc.php"

type AURInfoContainer struct {
	Type    string   `json:"type"`
	Results *AURInfo `json:"results"`
}

type AURInfo struct {
	Name           string
	Version        string
	Description    string
	ID             int64 `json:",string"`
	License        string
	NumVotes       int `json:",string"`
	OutOfDate      int `json:",string"`
	URL, URLPath   string
	FirstSubmitted int64 `json:",string"`
	LastModified   int64 `json:",string"`
}

type AURSearchContainer struct {
	Type    string         `json:"type"`
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	Name string
	ID   int64 `json:",string"`
}

func genericQuery(querytype, arg string, target interface{}) os.Error {
	form := make(url.Values)
	form.Add("type", querytype)
	form.Add("arg", arg)
	url := AURLocation + "?" + form.Encode()
	resp, er := http.Get(url)
	if er != nil {
		return er
	}

	reader := json.NewDecoder(resp.Body)
	reader.Decode(target)
	return nil
}

func GetInfo(pkg string) (*AURInfo, os.Error) {
  var info AURInfoContainer
  er := genericQuery("info", pkg, &info)
  return info.Results, er
}

func DoSearch(pattern string) ([]SearchResult, os.Error) {
  var info AURSearchContainer
  er := genericQuery("search", pattern, &info)
  return info.Results, er
}
