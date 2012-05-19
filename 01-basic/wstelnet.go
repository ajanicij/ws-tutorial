package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: wstelnet port")
		os.Exit(0)
	}
	
	port := os.Args[1]
	fmt.Println("Serving web on port", port)
	service := ":" + port
	
	http.Handle("/script/", http.FileServer(http.Dir(".")))
	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/", http.FileServer(http.Dir("./html/")))
	err := http.ListenAndServe(service, nil)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

