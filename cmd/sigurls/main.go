package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/drsigned/sigurls/pkg/runner"
	"github.com/logrusorgru/aurora/v3"
)

type options struct {
	listSources bool
	noColor     bool
	silent      bool
}

var (
	co options
	au aurora.Aurora
	so runner.Options
)

func banner() {
	fmt.Fprintln(os.Stderr, aurora.BrightBlue(`
     _                  _
 ___(_) __ _ _   _ _ __| |___
/ __| |/ _`+"`"+` | | | | '__| / __|
\__ \ | (_| | |_| | |  | \__ \
|___/_|\__, |\__,_|_|  |_|___/ v1.2.0
       |___/
`).Bold())
}

func init() {
	flag.StringVar(&so.Domain, "d", "", "")
	flag.StringVar(&so.Exclude, "exclude", "", "")
	flag.BoolVar(&co.listSources, "ls", false, "")
	flag.BoolVar(&co.noColor, "nc", false, "")
	flag.BoolVar(&co.silent, "silent", false, "")
	flag.BoolVar(&so.IncludeSubs, "subs", false, "")
	flag.StringVar(&so.Use, "use", "", "")

	flag.Usage = func() {
		banner()

		h := "USAGE:\n"
		h += "  sigurls [OPTIONS]\n"

		h += "\nOPTIONS:\n"
		h += "  -d             domain to fetch urls for\n"
		h += "  -exclude       comma(,) separated list of sources to exclude\n"
		h += "  -ls            list all the available sources\n"
		h += "  -nc            no color mode\n"
		h += "  -silent        silent mode: output urls only\n"
		h += "  -subs          include subdomains' urls\n"
		h += "  -use           comma(,) separated list of sources to use\n"

		fmt.Fprintf(os.Stderr, h)
	}

	flag.Parse()

	au = aurora.NewAurora(!co.noColor)
}

func main() {
	options, err := runner.ParseOptions(&so)
	if err != nil {
		log.Fatalln(err)
	}

	if !co.silent {
		banner()
	}

	if co.listSources {
		fmt.Println("[", au.BrightBlue("INF"), "] current list of the available", au.Underline(strconv.Itoa(len(options.YAMLConfig.Sources))+" sources").Bold())
		fmt.Println("[", au.BrightBlue("INF"), "] sources marked with an * needs key or token")
		fmt.Println("")

		keys := options.YAMLConfig.GetKeys()
		needsKey := make(map[string]interface{})
		keysElem := reflect.ValueOf(&keys).Elem()

		for i := 0; i < keysElem.NumField(); i++ {
			needsKey[strings.ToLower(keysElem.Type().Field(i).Name)] = keysElem.Field(i).Interface()
		}

		for _, source := range options.YAMLConfig.Sources {
			if _, ok := needsKey[source]; ok {
				fmt.Println(">", source, "*")
			} else {
				fmt.Println(">", source)
			}
		}

		os.Exit(0)
	}

	if !co.silent {
		fmt.Println("[", au.BrightBlue("INF"), "] fetching urls for", au.Underline(options.Domain).Bold())

		if options.IncludeSubs {
			fmt.Println("[", au.BrightBlue("INF"), "] -subs used: includes subdomains' urls")
		}

		fmt.Println("")
	}

	runner := runner.New(options)

	URLs, err := runner.Run()
	if err != nil {
		log.Fatalln(err)
	}

	for URL := range URLs {
		if co.silent {
			fmt.Println(URL)
		} else {
			fmt.Println(fmt.Sprintf("[%s] %s", au.BrightBlue(URL.Source), URL.Value))
		}
	}
}
