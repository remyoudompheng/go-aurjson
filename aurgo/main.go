// aurgo is a small AUR standalone client.
package main

import (
	aurjson "github.com/remyoudompheng/go-aurjson"
	"flag"
	"fmt"
	"os"
	"text/template"
)

const infoTplString = `Name          : {{ .Name }}
Version       : {{ .Version }}
Description   : {{ .Description }}
License       : {{ .License }}
URL           : https://aur.archlinux.org{{ .URLPath }}
Upstream      : {{ .URL }}
Maintainer    : {{ .Maintainer }}
Last Modified : {{ .LastModified }}
`

var infoTpl *template.Template

func init() {
	var er error
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
			fmt.Printf("Error: %s\n", er)
			return
		}
		for _, item := range results {
			fmt.Printf("%s %s\n", item.Name, item.Version)
			fmt.Printf("  %s\n", item.Description)
		}
	case infoarg != "":
		info, er := aurjson.GetInfo(infoarg)
		if er != nil {
			fmt.Printf("Error: %s\n", er)
			return
		}
		infoTpl.Execute(os.Stdout, *info)
	default:
		flag.Usage()
	}
}
