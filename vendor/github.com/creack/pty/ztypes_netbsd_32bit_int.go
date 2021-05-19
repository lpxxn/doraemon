// +build netbsd
// +build 386 amd64 arm arm64

package pty

type ptmget struct {
	Cfd int32
	Sfd int32
	Cn  [1024]int8
	Sn  [1024]int8
}

var ioctl_TIOCPTSNAME = 0x48087448
var ioctl_TIOCGRANTPT = 0x20007447
