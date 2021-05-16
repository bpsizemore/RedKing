package main

import(
    "net/http"
	"fmt"
	"strconv"
	"flag"
//	"os"
	"errors"
	"log"
	"strings"
)

var Config struct {
	asciiart string
	redirectStatus int
	verbosity int
	redirectUrl string
	redirectHost string
	port int
	httpPort string
	mode string
	validModes []string
	scanPorts []string
}

func ValidateUrl(urlPtr *string) error {
	if (*urlPtr == "") {
		return errors.New("No URL supplied.")
	}
	return nil
}

func ValidateMode(modePtr *string) error {
	for _, mode := range Config.validModes {
		if strings.ToLower(*modePtr) == strings.ToLower(mode) {
			return nil 
		}
	}
	return errors.New("Invalid mode selected. Please choose one of the valid modes")// + string(Config.validModes))
}

func ExtractHostFromUrl(urlPtr *string) string {
	return *urlPtr
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are being redirected...")
}

func RedirectToSite(url string) {
	// Static Redirect to a single site.
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, Config.redirectStatus)
		log.Printf("Request from %s", r.RemoteAddr)
	})
}

func RedirectToSiteWithTarget(target string,url string) {
	// Static Redirect to a single site.
	http.HandleFunc(target, func (w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, Config.redirectStatus)
		fmt.Printf("Redirection from %s to %s\n", r.RemoteAddr, url)
	})
}

func RedirectPortScan(url string) {
	// Used to do a simple "port scan" against a host via redirect.
	// Use this in conjunction with Burp or an automated script for fastest results.
	// Defaults to testing a small number of important ports.

	for i, port := range Config.scanPorts {
		if Config.verbosity > 0 {
			fmt.Println("adding handler for port ", port, i)
		}
		//Setup a bunch of handlers like /1, /2, /3 which each redirect to a different port on one host
		RedirectToSiteWithTarget("/" + strconv.Itoa(i) , url + ":" + port)
	}

}

func main() {

	// Setup default config - not sure where this should go instead
	Config.scanPorts = []string{"22", "80", "443", "445", "3389", "8000", "8080"}
	Config.validModes = []string{"single", "portscan"}
	Config.asciiart = 
`

______         _   _   ___             
| ___ \       | | | | / (_)            
| |_/ /___  __| | | |/ / _ _ __   __ _ 
|    // _ \/ _' | |    \| | '_ \ / _' |
| |\ \  __/ (_| | | |\  \ | | | | (_| |
\_| \_\___|\__,_| \_| \_/_|_| |_|\__, |
                                  __/ |
                                 |___/ 

`

	redirectStatusPtr := flag.Int("r", 302, "Redirect status code - suggested 301, 302, or 307")
	verbosityPtr := flag.Bool("v", false, "Verbose")
	urlPtr := flag.String("url", "",  "The URL used for redirects")
	portPtr := flag.Int("p", 8080, "The port used to host the redirect server")
	modePtr := flag.String("mode", "single", "The mode RedKing should execute in. Select from:\nsingle - redirect to a single URL\nportscan - create a series of redirects at localhost/1,localhost/2,...\n\tEach number will redirect to a different port on the target host.\nThe built in port scan ports are: 22,80,443,445,3389,8000,8080\n")

	flag.Parse()

	err := ValidateUrl(urlPtr)
	if (err != nil) {
		log.Fatal(err)
	}
	Config.redirectUrl = *urlPtr
	Config.redirectHost = ExtractHostFromUrl(urlPtr)

	Config.redirectStatus = *redirectStatusPtr
	if *verbosityPtr == false {
		Config.verbosity = 0
	} else {
		Config.verbosity = 1
	}
	
	Config.port = *portPtr
	Config.httpPort = ":" + strconv.Itoa(*portPtr)
	
	mode_err := ValidateMode(modePtr)
	if mode_err != nil {
		log.Fatal(mode_err)
	}
	Config.mode = *modePtr

	// 3 redirection types, 301, 302, and 307
	// default to 302

	fmt.Println(Config.asciiart)
	fmt.Printf("Mode: %s \nURL: %s\nPort: %s\n\n",Config.mode, Config.redirectUrl, Config.httpPort)

	switch Config.mode {
	case "single":
		RedirectToSite(Config.redirectUrl)
	case "portscan":
		RedirectPortScan(Config.redirectUrl)
	}

	fmt.Printf("Starting server on localhost%s\n",Config.httpPort)
	http.ListenAndServe(Config.httpPort, nil)

}