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

	SendMsg(false, "", fmt.Sprintf("Serving %s on HTTP port: %d\n", folderPath, listener.Addr().(*net.TCPAddr).Port), Cyan, false)
	addr := fmt.Sprintf("http://%s:%d", ip.String(), listener.Addr().(*net.TCPAddr).Port)
	SendMsg(false, "open in brower, or scan QR", addr, Yellow, false)
	RenderQRString(addr)
	return http.Serve(listener, nil)
}

func PrivateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
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
