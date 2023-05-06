package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GIPInfo struct {
	Ip      string // IP地址
	Address string // 地址(物理)
	ISP     string // 运营商
}

func (gi GIPInfo) String() string {
	return fmt.Sprintf("\n公网IP: %s\n  地址: %s\n运营商: %s\n", gi.Ip, gi.Address, gi.ISP)
}

type GIPSrv struct {
}

func (gs *GIPSrv) GetPublicIpByHttp() (info *GIPInfo, err error) {
	req, _ := http.NewRequest(http.MethodGet, "https://cip.cc", nil)
	// 有趣, 对方检查了代理
	req.Header.Add("User-Agent", "curl/*")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respStr := string(respBytes)
	info = new(GIPInfo)
	cntMatched := 0
	lineSlice := strings.Split(respStr, "\n")
	for _, line := range lineSlice {
		kv := strings.Split(strings.TrimSpace(line), ":")
		if len(kv) == 2 {
			k := strings.TrimSpace(kv[0])
			v := strings.TrimSpace(kv[1])
			switch k {
			case "IP":
				info.Ip = v
				cntMatched++
			case "地址":
				info.Address = v
				cntMatched++
			case "运营商":
				info.ISP = v
				cntMatched++
			}
		}
	}

	if cntMatched == 0 {
		info = nil
		err = errors.New("查询失败")
		return
	}
	return
}
