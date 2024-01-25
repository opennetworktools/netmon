package internal

import (
	"fmt"
	"log"
	"net"
	"path/filepath"

	"github.com/oschwald/maxminddb-golang"
)

func GetASN(IPAddress string) {
	mmdbASNReader(IPAddress)
}

func GetCountry(IPAddress string) {
	mmdbCountryReader(IPAddress)
}

func mmdbASNReader(IPAddress string) {
	dbFilePath := filepath.Join("mmdb", "GeoLite2-ASN.mmdb")
	db, err := maxminddb.Open(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(IPAddress)

	var record struct {
		ASN uint  `maxminddb:"autonomous_system_number"`
		ASName string `maxminddb:"autonomous_system_organization"`
	}
	err = db.Lookup(ip, &record)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%d %s \n", record.ASN, record.ASName)
}

func mmdbCountryReader(IPAddress string) {
	// rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    // if err != nil {
    //     log.Fatal(err)
    // }
    // dbFilePath := filepath.Join(rootDir, "mmdb", "GeoLite2-Country.mmdb")
    // file, err := os.Open(dbFilePath)
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer file.Close()
	
	dbFilePath := filepath.Join("mmdb", "GeoLite2-Country.mmdb")
	db, err := maxminddb.Open(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(IPAddress)

	var record struct {
		Country struct {
			Names map[string]string `maxminddb:"names"`
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	err = db.Lookup(ip, &record)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%s %s \n", record.Country.ISOCode, record.Country.Names["en"])
}
