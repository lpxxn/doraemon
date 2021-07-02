package main

import (
	"flag"

	"github.com/lpxxn/doraemon/internal"
)

func main() {
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	if err := internal.HttpFileServ(*directory); err != nil {
		panic(err)
	}
	//ip, err := internal.PrivateIPv4()
	//if err != nil {
	//	panic(err)
	//}
	//
	//listener, err := net.Listen("tcp", ":0")
	//if err != nil {
	//	panic(err)
	//}
	//http.Handle("/", http.FileServer(http.Dir(*directory)))
	//
	//log.Printf("Serving %s on HTTP port: %d\n", *directory, listener.Addr().(*net.TCPAddr).Port)
	//addr := fmt.Sprintf("http://%s:%d", ip.String(), listener.Addr().(*net.TCPAddr).Port)
	//fmt.Println(addr)
	//log.Fatal(http.Serve(listener, nil))
}
