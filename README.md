# sigurls

![made with go](https://img.shields.io/badge/made%20with-Go-0040ff.svg) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg) [![open issues](https://img.shields.io/github/issues-raw/drsigned/sigurls.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigurls/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/drsigned/sigurls.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigurls/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?colorB=0040FF)](https://github.com/drsigned/sigurls/blob/master/LICENSE) [![twitter](https://img.shields.io/badge/twitter-@drsigned-0040ff.svg)](https://twitter.com/drsigned)

sigurls fetches known URLs from **AlienVault's OTX**, **Common Crawl**, **URLScan** and the **Wayback Machine** for any given domain.

## Resources

* [Installation](#installation)
    * [From Binary](#from-binary)
    * [From source](#from-source)
    * [From github](#from-github)
* [Usage](#usage)
* [Contribution](#contribution)

## Installation

#### From Binary

You can download the pre-built binary for your platform from this repository's [releases](https://github.com/drsigned/sigurls/releases/) page, extract, then move it to your `$PATH`and you're ready to go.

#### From Source

sigurls requires **go1.14+** to install successfully. Run the following command to get the repo

```bash
$ GO111MODULE=on go get -u -v github.com/drsigned/sigurls/cmd/sigurls
```

#### From Github

```bash
$ git clone https://github.com/drsigned/sigurls.git; cd sigurls/cmd/sigurls/; go build; mv sigurls /usr/local/bin/; sigurls -h
```

## Usage

To display help message for sigurls use the `-h` flag:

```
$ sigurls -h

     _                  _
 ___(_) __ _ _   _ _ __| |___
/ __| |/ _` | | | | '__| / __|
\__ \ | (_| | |_| | |  | \__ \
|___/_|\__, |\__,_|_|  |_|___/ v1.0.0
       |___/

USAGE:
  sigurls [OPTIONS]

OPTIONS:
  -d              domain to fetch urls for
  -e              comma separated list of sources to exclude
  -nc             no color mode
  -s              silent mode: output urls only
  -subs           include subdomains' urls
  -u              comma separated list of sources to use
```

## Contribution

[Issues](https://github.com/drsigned/sigurls/issues) and [Pull Requests](https://github.com/drsigned/sigurls/pulls) are welcome!
