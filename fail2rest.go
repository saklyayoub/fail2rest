package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Sean-Der/fail2go"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type Configuration struct {
	Addr           string
	Fail2banSocket string
}

type ErrorBody struct {
	Error string
}

var fail2goConn *fail2go.Conn

func main() {
	configPath := flag.String("config", "config.json", "path to config.json")
	flag.Parse()

	file, fileErr := os.Open(*configPath)

	if fileErr != nil {
		fmt.Println("failed to open config:", fileErr)
		os.Exit(1)
	}

	configuration := new(Configuration)
	configErr := json.NewDecoder(file).Decode(configuration)

	if configErr != nil {
		fmt.Println("config error:", configErr)
		os.Exit(1)
	}

	fail2goConn := fail2go.Newfail2goConn(configuration.Fail2banSocket)
	r := mux.NewRouter()

	globalHandler(r.PathPrefix("/global").Subrouter(), fail2goConn)
	jailHandler(r.PathPrefix("/jail").Subrouter(), fail2goConn)

	http.Handle("/", r)
	fmt.Println(http.ListenAndServe(configuration.Addr, nil))
}
