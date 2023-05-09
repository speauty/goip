package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

const (
	Protocol   = "udp4"
	NameServer = "ns1.google.com:53"
	IPServer   = "o-o.myaddr.l.google.com"
)

type GIPInfo struct {
	LocalIp string // 本地IP
	Ip      string // IP地址
	Address string // 地址(物理)
	ISP     string // 运营商
}

func (gi GIPInfo) GetLocalIp() string {
	if gi.LocalIp == "" {
		return "暂无"
	}
	return gi.LocalIp
}

func (gi GIPInfo) GetPublicIp() string {
	if gi.Ip == "" {
		return "暂无"
	}
	return gi.Ip
}

func (gi GIPInfo) GetAddress() string {
	if gi.Address == "" {
		return "暂无"
	}
	return gi.Address
}

func (gi GIPInfo) GetISP() string {
	if gi.ISP == "" {
		return "暂无"
	}
	return gi.ISP
}

func (gi GIPInfo) String() string {
	return fmt.Sprintf(
		"\n本地IP: %s\n公网IP: %s\n  地址: %s\n运营商: %s\n",
		gi.GetLocalIp(), gi.GetPublicIp(), gi.GetAddress(), gi.GetISP(),
	)
}

type GIPSrv struct {
}

// GetPublicIpByHttp 通过模拟curl请求获取公网IP(不推荐)
func (gs *GIPSrv) GetPublicIpByHttp() (info *GIPInfo, err error) {
	req, _ := http.NewRequest(http.MethodGet, "https://cip.cc", nil)
	// 有趣, 对方检查了代理(关键)
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

	if localIp, _ := gs.GetLocalIp(); localIp != "" {
		info.LocalIp = localIp
	}
	return
}

// GetPublicIpByDial
// 参考: https://github.com/ysmood/myip
func (gs *GIPSrv) GetPublicIpByDial() (info *GIPInfo, err error) {
	r := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, Protocol, NameServer)
		},
	}
	txt, err := r.LookupTXT(context.Background(), IPServer)
	if err != nil {
		return
	}

	if len(txt) == 0 {
		err = errors.New("获取公网IP失败")
		return
	}
	info = new(GIPInfo)
	info.Ip = txt[0]
	if localIp, _ := gs.GetLocalIp(); localIp != "" {
		info.LocalIp = localIp
	}
	return
}

func (gs *GIPSrv) GetLocalIp() (ip string, err error) {
	conn, err := net.Dial(Protocol, NameServer)
	if err != nil {
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = localAddr.IP.String()
	return
}
