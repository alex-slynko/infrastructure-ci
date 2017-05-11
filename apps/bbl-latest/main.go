package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry/infrastructure-ci/apps/bbl-latest/utils"
)

func main() {
	port := os.Getenv("PORT")
	fmt.Printf("Starting server on port: %s...\n", port)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			latestBBLVersion string
			lastModified     string
		)

		switch req.URL.Path {
		case "/latest":
			queryParams := req.URL.Query()
			os := queryParams.Get("os")

			var err error
			// NOT TESTED: saving and passing latestBBLVersion and lastModified
			latestBBLVersion, lastModified, err = utils.LatestBBLVersion(latestBBLVersion, lastModified)
			if err != nil {
				// not tested
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to query github.com for latest bbl version. Try again later."))
				fmt.Println(err.Error())
			}

			osTypes := map[string]string{
				"linux": "linux_x86-64",
				"osx":   "osx",
			}
			osType, ok := osTypes[os]
			if ok {
				redirectURL := fmt.Sprintf("https://github.com/cloudfoundry/bosh-bootloader/releases/download/%[1]s/bbl-%[1]s_%s", latestBBLVersion, osType)
				http.Redirect(w, req, redirectURL, http.StatusFound)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("missing required query parameter: os=[osx,linux]"))
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	if err != nil {
		log.Fatal(err)
	}
}
