# sigurls

![made with go](https://img.shields.io/badge/made%20with-Go-0040ff.svg) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg) [![open issues](https://img.shields.io/github/issues-raw/drsigned/sigurls.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigurls/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/drsigned/sigurls.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigurls/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?colorB=0040FF)](https://github.com/drsigned/sigurls/blob/master/LICENSE) [![twitter](https://img.shields.io/badge/twitter-@drsigned-0040ff.svg)](https://twitter.com/drsigned)

sigurls is a reconnaissance tool, it fetches URLs from **AlienVault's OTX**, **Common Crawl**, **URLScan**, **Github** and the **Wayback Machine**.

## Resources

* [Usage](#usage)
* [Installation](#installation)
    * [From Binary](#from-binary)
    * [From source](#from-source)
    * [From github](#from-github)
* [Post Installtion](#post-installation)
* [Contribution](#contribution)

## Usage

To display help message for sigurls use the `-h` flag:

```
$ sigurls -h

     _                  _
 ___(_) __ _ _   _ _ __| |___
/ __| |/ _` | | | | '__| / __|
\__ \ | (_| | |_| | |  | \__ \
|___/_|\__, |\__,_|_|  |_|___/ v1.3.0
       |___/

USAGE:
  sigurls [OPTIONS]

OPTIONS:
  -d             domain to fetch urls for
  -sE            comma(,) separated list of sources to exclude
  -iS            include subdomains' urls
  -sL            list all the available sources
  -nC            no color mode
  -silent        silent mode: output urls only
  -sU            comma(,) separated list of sources to use
```

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

## Post Installation

sigurls will work after [installation](#installation). However, to configure sigurls to work with certain services - currently github - you will need to have setup API keys. The API keys are stored in the `$HOME/.config/sigurls/conf.yaml` file - created upon first run - and uses the YAML format. Multiple API keys can be specified for each of these services.

Example:

```yaml
version: 1.3.0
sources:
    - commoncrawl
    - github
    - otx
    - urlscan
    - wayback
keys:
    github:
        - d23a554bbc1aabb208c9acfbd2dd41ce7fc9db39
        - asdsd54bbc1aabb208c9acfbd2dd41ce7fc9db39
```

## Contribution

[Issues](https://github.com/drsigned/sigurls/issues) and [Pull Requests](https://github.com/drsigned/sigurls/pulls) are welcome!
