package main

import (
	"fmt"
	"github.com/satori/go.uuid"
)
var countryMap map[string]string
var countrymap = make(map[string]string)

func main() {
	u1 := uuid.Must(uuid.NewV4()).String()

	countrymap["ShengRI"] = u1

	for country := range countryMap {
		fmt.Println("Capital of",country,"is",countrymap[country])
	}


}