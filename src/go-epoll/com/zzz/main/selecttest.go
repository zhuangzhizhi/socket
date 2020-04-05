package main

import (
	"fmt"
	"go-epoll/go-epoll/com/zzz/util"
	"net"
	"syscall"
	"time"
)

func main() {

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("syscall.Socket err=", err)
	}
	fmt.Println("server socket fd=", fd)
	syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	//if err = syscall.SetNonblock(fd, true); err != nil {
	//	fmt.Println("setnonblock1: ", err)
	//	os.Exit(1)
	//}
	BACKLOG := 1024
	addr := syscall.SockaddrInet4{
		Port: 8888,
	}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	syscall.Bind(fd, &addr)
	syscall.Listen(fd, BACKLOG)

	//for i:=0; i < BACKLOG; i++ {
	//	nfd, _, err := syscall.Accept(fd)
	//	if err != nil {
	//		fmt.Println("syscall.Accept err=", err)
	//		continue
	//	}
	//	acceptedFdSet[i] = nfd
	//	if nfd > maxFd {
	//		maxFd = nfd
	//	}
	//}
	masterSet := syscall.FdSet{}
	//util.FD_ZERO(&masterSet)
	for i := range masterSet.Bits {
		masterSet.Bits[i] = 0xffffffff
	}
	util.FD_SET(fd, &masterSet)
	var maxFd = fd
	//timeout := syscall.Timeval{Sec: 2, Usec: 0}
	for {
		var workingSet = syscall.FdSet{Bits: masterSet.Bits}
		fmt.Println("before select", workingSet)
		n, err := syscall.Select(maxFd + 1, &workingSet, nil, nil, nil)
		fmt.Println("after select", workingSet)
		time.Sleep(time.Second * 2)
		if err != nil {
			fmt.Println("syscall.Select, err=", err)
			continue
		}
		if n < 0 {
			fmt.Println("syscall.Select，, n=", n)
			break
		}


	}
}




//func FD_ISSET(fd int, p *syscall.FdSet) bool {
//	return (p.Bits[fd/64] & (1 << (uint(fd) % 64))) != 0
//}
//
////func FD_ISSET(sysfd int, set *syscall.FdSet) bool {
////	s := unsafe.Pointer(set)
////	fd := C.int(sysfd)
////	return C._FD_ISSET(fd, s) != 0
////}
//
//func FD_SET(sysfd int, set *syscall.FdSet) {
//	var shang int64 = int64(sysfd / 64)
//	var yushu int64 = int64(sysfd % 64)
//	set.Bits[shang] = set.Bits[shang] | yushu
//}
//
//func FD_ZERO(set *syscall.FdSet) {
//	var v [16]int64
//	set.Bits = v
//}

//func FD_SET(fd int, p *syscall.FdSet) {
//	fmt.Println("FD_SET, fd=", fd, " rfset=", p)
//	p.Bits[fd/64] = 1 << (fd % 64)
//	fmt.Println("FD_SET, fd=", fd, " rfset=", p)
//}
//
//func FD_ISSET(i int, p *syscall.FdSet) bool {
//	return (p.Bits[i/64] & (1 << uint(i) % 64)) != 0
//}
//
//func FD_ZERO(p *syscall.FdSet) {
//	for i := range p.Bits {
//		p.Bits[i] = 0
//	}
//}


