package main

import (
	"os"
	"fmt"
	"flag"
	"template"
	aurjson "archlinux/aurjson"
)

const infoTplString = `
Name     : {{ .Name }}
Version  : {{ .Version }}
URL      : https://aur.archlinux.org{{ .URLPath }}
Upstream : {{ .URL }}
`

var infoTpl *template.Template

func init() {
	var er os.Error
	infoTpl, er = template.New("info").Parse(infoTplString)
	if er != nil {
		fmt.Printf("couldn't compile template: %s\n", er)
		panic("template parsing error")
	}
}

func main() {
	var searchstr, infoarg string
	flag.StringVar(&searchstr, "s", "", "search packages by pattern")
	flag.StringVar(&infoarg, "i", "", "get information for a specific package")
	flag.Parse()

	switch {
	case searchstr != "":
		results, er := aurjson.DoSearch(searchstr)
		if er != nil {
			fmt.Printf("Error: %s", er)
			return
		}
		for _, item := range results {
			fmt.Printf("%s\n", item.Name)
		}
	case infoarg != "":
		info, er := aurjson.GetInfo(infoarg)
		if er != nil {
			fmt.Printf("Error: %s", er)
			return
		}
		infoTpl.Execute(os.Stdout, *info)
	default:
		flag.Usage()
	}
}
