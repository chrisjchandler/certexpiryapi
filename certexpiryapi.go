package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the URL from the request parameters
		url := r.URL.Query().Get("url")

		// Check the certificate expiration for the given URL
		daysRemaining, err := checkCertExpiration(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the number of days remaining as a response
		fmt.Fprintf(w, "%d", daysRemaining)
	})

	http.ListenAndServe(":8080", nil)
}

func checkCertExpiration(url string) (int, error) {
	// Make an HTTPS request to the given URL
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Get the certificate from the response
	cert := resp.TLS.PeerCertificates[0]

	// Calculate the number of days remaining before the certificate expires
	now := time.Now()
	duration := cert.NotAfter.Sub(now)
	daysRemaining := int(duration.Hours() / 24)

	return daysRemaining, nil
}
