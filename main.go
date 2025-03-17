package main

import (
	"fmt"
	"github.com/likexian/whois"
	"github.com/oschwald/geoip2-golang"
	"net"
	"net/http"
)

func getClientIP(r *http.Request) string {
	// Try to get the IP from the X-Forwarded-For header (useful if behind a proxy)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}

	// Try to get the IP from the X-Real-IP header (useful if behind a proxy)
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fallback to the remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	fmt.Fprintf(w, "%s\n", clientIP)
}

func detailsHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)

	// Get whois information
	whoisInfo, err := whois.Whois(clientIP)
	if err != nil {
		whoisInfo = fmt.Sprintf("Error getting whois information: %v", err)
	}

	// Get GeoIP information
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		fmt.Fprintf(w, "Error opening GeoIP database: %v\n", err)
		return
	}
	defer db.Close()

	ip := net.ParseIP(clientIP)
	record, err := db.City(ip)
	if err != nil {
		fmt.Fprintf(w, "Error getting GeoIP information: %v\n", err)
		return
	}

	// Get DNS reverse lookup information
	names, err := net.LookupAddr(clientIP)
	if err != nil {
		fmt.Fprintf(w, "Error performing DNS reverse lookup: %v\n", err)
		return
	}

	fmt.Fprintf(w, "Your IP address is: %s\n", clientIP)

	// Print DNS reverse lookup information
	fmt.Fprintf(w, "\nDNS Reverse Lookup:\n")
	for _, name := range names {
		fmt.Fprintf(w, "%s\n", name)
	}

	// Print client request headers
	fmt.Fprintf(w, "\nClient Request Headers:\n")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\n", name, value)
		}
	}

	// Print whois information
	fmt.Fprintf(w, "\nWhois Information:\n%s\n", whoisInfo)

	// Print GeoIP information
	fmt.Fprintf(w, "\nGeoIP Information:\n")
	fmt.Fprintf(w, "Country: %s\n", record.Country.Names["en"])
	fmt.Fprintf(w, "City: %s\n", record.City.Names["en"])
	fmt.Fprintf(w, "Latitude: %f\n", record.Location.Latitude)
	fmt.Fprintf(w, "Longitude: %f\n", record.Location.Longitude)
}

func main() {
	http.HandleFunc("/", ipHandler)
	http.HandleFunc("/details", detailsHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
