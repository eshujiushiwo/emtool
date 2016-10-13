package main

import (
	"emtool/mlogger"
	"flag"
	"fmt"
)

func main() {
	var fromhost, tohost, fromport, toport string
	flag.StringVar(&fromhost, "fromhost", "", "the source host")
	flag.StringVar(&tohost, "tohost", "", "the dest host")
	flag.StringVar(&fromport, "fromport", "27017", "the source port")
	flag.StringVar(&toport, "toport", "27017", "the dest port")

	flag.Parse()

	if fromhost == "" || tohost == "" {
		fmt.Println("Please use -help to check the usage")
	}

}
