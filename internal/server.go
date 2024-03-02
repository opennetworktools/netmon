package internal

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Archived
func StartServer(c chan CPacket, m map[string]CHost) {
	pagePath := filepath.Join("pages", "index.html")
	printPackets := func(wr http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(pagePath))
		tmpl.Execute(wr, c)
	}
	http.HandleFunc("/", printPackets)

	pagePath1 := filepath.Join("pages", "hosts.html")
	printHosts := func(wr http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(pagePath1))
		tmpl.Execute(wr, m)
	}
	http.HandleFunc("/hosts", printHosts)
	
	fmt.Println("Page is live on PORT 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}