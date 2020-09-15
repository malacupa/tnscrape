# tnscrape

TLS certificate Subject Name and Subject Alternative Name scraper. Supply list of IP and ports and receive list of either domains or all other possible data stored in TLS certificates. If you also supply list of domains you might get more data.

Works on TLSv1.0 - TLSv1.3 due to Golang's TLS support. See files in this repo for example `ipp.csv` and `ipd.csv`.

## Install

```bash
go get github.com/malacupa/tnscrape
```

## Usage

```bash
./tnscrape
Usage: ./tnscrape [options] <ipp.csv> [<ipd.csv>]
  -a	include all information from certificates (IPs, emails, URLs), not just domains
<ipp.csv>	File in CSV format having first column IP and next ports to scrape on the IP
<ipd.csv>	File in CSV format having first column IP and next domains associated with the IP
```

Examples:
```
$ ./tnscrape ipp-example.csv
invalid2.invalid
www.github.com
github.com
github.io
githubusercontent.com
htbridge.ch
www.htbridge.ch
```

```
$ ./tnscrape ipp-example.csv ipd-example.csv
invalid2.invalid
misc-sni.google.com
1ucrs.com
abc.xyz
adsensecustomsearchads.com
advertisercommunity.com
...continues...
datafusion-api-staging.cloud.google.com
datafusion-api.cloud.google.com
go-lang.net
go-lang.org
golang.com
golang.google.cn
golang.net
golang.org
google-syndication.com
www.github.com
github.com
github.io
githubusercontent.com
yourbasic.org
htbridge.ch
www.htbridge.ch
immuniweb.com
```

## Similar Tools

 * [HostHunter](https://github.com/SpiderLabs/HostHunter) from SipderLabs
 * [certasset](https://github.com/arbazkiraak/certasset) from arbazkiraak
 * [sslScrape](https://github.com/cheetz/sslScrape) from cheetz
 * ...

## License
[MIT](https://choosealicense.com/licenses/mit/)

