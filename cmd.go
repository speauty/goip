package main

import (
	"fmt"
	"log"
)

func main() {
	res, err := new(GIPSrv).GetPublicIpByHttp()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	return
}
