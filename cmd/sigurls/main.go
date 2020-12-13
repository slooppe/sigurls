package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/drsigned/sigurls/pkg/sigurls"
	"github.com/logrusorgru/aurora/v3"
)

type options struct {
	noColor bool
	silent  bool
}

var (
	co options
	au aurora.Aurora
	so sigurls.Options
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
	flag.StringVar(&so.ExcludeSources, "e", "", "")
	flag.BoolVar(&co.noColor, "nc", false, "")
	flag.BoolVar(&co.silent, "s", false, "")
	flag.BoolVar(&so.IncludeSubs, "subs", false, "")
	flag.StringVar(&so.UseSources, "u", "", "")

	flag.Usage = func() {
		banner()

		h := "USAGE:\n"
		h += "  sigurls [OPTIONS]\n"

		h += "\nOPTIONS:\n"
		h += "  -d              domain to fetch urls for\n"
		h += "  -e              comma separated list of sources to exclude\n"
		h += "  -nc              no color mode\n"
		h += "  -s              silent mode: output urls only\n"
		h += "  -subs           include subdomains' urls\n"
		h += "  -u              comma separated list of sources to use\n"

		fmt.Fprintf(os.Stderr, h)
	}

	flag.Parse()

	au = aurora.NewAurora(!co.noColor)
}

func main() {
	options, err := sigurls.ParseOptions(&so)
	if err != nil {
		log.Fatalln(err)
	}

	if !co.silent {
		banner()

		fmt.Println("[", au.BrightBlue("INF"), "] fetching urls for", au.Underline(options.Domain).Bold())

		if options.IncludeSubs {
			fmt.Println("[", au.BrightBlue("INF"), "] -subs used: includes subdomains' urls")
		}

		fmt.Println("")
	}

	runner := sigurls.NewRunner(options)

	URLs, err := runner.Run()
	if err != nil {
		log.Fatalln(err)
	}

	for n := range URLs {
		if co.silent {
			fmt.Println(n.URL)
		} else {
			fmt.Println(fmt.Sprintf("[%s] %s", au.BrightBlue(n.Source), n.URL))
		}
	}
}
