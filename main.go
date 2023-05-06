package main

import (
	"fmt"
)

func main() {
	res, err := new(GIPSrv).GetPublicIpByHttp()
	if err != nil {
		fmt.Println("查询异常, 错误:", err)
		return
	}
	fmt.Println(res)
	return
}
