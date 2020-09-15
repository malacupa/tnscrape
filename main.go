package main

import (
	"fmt"
	"net"
	"crypto/tls"
	"crypto/x509"
	"os"
	"strings"
	"sync"
	"regexp"
	"flag"
	"io/ioutil"
)

func delWildcard(str string) string {
	wcRe := regexp.MustCompile(`^\*\.`)
	return wcRe.ReplaceAllString(str, "")
}

func stringInSlice(a string, list []string) bool {
  // kudos https://stackoverflow.com/questions/15323767/does-go-have-if-x-in-construct-similar-to-python
  for _, b := range list {
    if b == a {
      return true
    }
  }
  return false
}

func myUsage() {
	fmt.Printf("Usage: %s [options] <ipp.csv> [<ipd.csv>]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("<ipp.csv>\tFile in CSV format having first column IP and next ports to scrape on the IP")
	fmt.Println("<ipd.csv>\tFile in CSV format having first column IP and next domains associated with the IP")
}

// TODO docker TLS testbed
func main() {

	flag.Usage = myUsage
	allInfoPtr := flag.Bool("a", false, "include all information from certificates (IPs, emails, URLs), not just domains")
	flag.Parse()

	if flag.NArg() < 1 || flag.NArg() > 2 {
    flag.Usage()
    os.Exit(1)
  }

	ipps, err := ioutil.ReadFile(flag.Arg(0))
  if err != nil {
    panic(err)
  }

	var ipds []byte
	if flag.NArg() > 1 {
		var err error
		ipds, err = ioutil.ReadFile(flag.Arg(1))
		if err != nil {
		  panic(err)
		}
	} else {
		ipds = nil
	}

	tlsConf := &tls.Config {
		InsecureSkipVerify: true,
	}
	dialer := &net.Dialer {
		Timeout: 2500000000,	// 2.5s
	}

	var wg sync.WaitGroup
	out := make(chan []string)
	addr := make(chan string)

	// TODO maybe do flag for concurrency
	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func () {
			defer wg.Done()

			for a := range addr {
				func () {
					conn, err := tls.DialWithDialer(dialer, "tcp", a, tlsConf)
					if err != nil {
						// if TLS cant be established this throws error
						return
					}

					defer conn.Close()

					var cert x509.Certificate
					var o []string

					state := conn.ConnectionState()
					// leaf certificate is first
					cert = *state.PeerCertificates[0]

					o = append(o, delWildcard(cert.Subject.CommonName))
					for _, dn := range cert.DNSNames {
						o = append(o, delWildcard(dn))
					}

					if *allInfoPtr {
						for _, email := range cert.EmailAddresses {
							o = append(o, email)
						}
						for _, ip := range cert.IPAddresses {
							o = append(o, ip.String())
						}
						for _, uri := range cert.URIs {
							o = append(o, uri.String())
						}
					}

					out <- o
				}()
			}
		}()
	}

	go func() {
		ippLines := strings.Split(string(ipps), "\n")
		ipdMap := make(map[string][]string)

		if ipds != nil {
			ipdLines := strings.Split(string(ipds), "\n")
			for _, line := range ipdLines {
				lineArr := strings.Split(line, ",")
				ip := lineArr[0]
				ipdMap[ip] = lineArr[1:]
			}
		}

		for _, line := range ippLines {
			lineArr := strings.Split(line, ",")
			ip := lineArr[0]
			for _, port := range lineArr[1:] {
				addr <- ip+":"+port
				for _, host := range ipdMap[ip] {
					addr <- host+":"+port
				}
			}
		}
		close(addr)
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	var items []string
	for o := range out {
		for _, item := range o {
			if ! stringInSlice(item, items) {
				items = append(items, item)
			}
		}
	}

	for _, item := range items {
		fmt.Println(item)
	}
}

