package internal

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

func HttpFileServ(folderPath string) error {
	ip, err := PrivateIPv4()
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	http.Handle("/", http.FileServer(http.Dir(folderPath)))

	addr := fmt.Sprintf("http://%s:%d", ip.String(), listener.Addr().(*net.TCPAddr).Port)
	SendMsg(false, "open in browser, or scan QR ", addr, Yellow, true)
	RenderQRString(addr)
	return http.Serve(listener, nil)
}

func PrivateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipNet, ok := a.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}

		ip := ipNet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}
