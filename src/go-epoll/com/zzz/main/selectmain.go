package main

import (
	"fmt"
	"go-epoll/go-epoll/com/zzz/util"
	"net"
	"syscall"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("syscall.Socket err=", err)
	}

	//if err = syscall.SetNonblock(fd, true); err != nil {
	//	fmt.Println("setnonblock1: ", err)
	//	os.Exit(1)
	//}
	totalNum := 1
	addr := syscall.SockaddrInet4{Port: 8888}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	syscall.Bind(fd, &addr)
	syscall.Listen(fd, totalNum)
	fdSet := make([]int, 1024)
	maxFd := 0
	for i:=0; i < totalNum; i++ {
		nfd, _, err := syscall.Accept(fd)
		if err != nil {
			fmt.Println("syscall.Accept err=", err)
			continue
		}
		fdSet[i] = nfd
		if nfd > maxFd {
			maxFd = nfd
		}
	}
	rfdSet := syscall.FdSet{}
	for {
		util.FD_ZERO(&rfdSet)
		for i:=0; i < totalNum; i++ {
			util.FD_SET(fdSet[i], &rfdSet)
		}
		_, err = syscall.Select(maxFd+1, &rfdSet, nil, nil, nil)
		for i:=0; i < totalNum; i++ {
			if util.FD_ISSET(fdSet[i], &rfdSet) {
				fmt.Println("echo...")
				go util.ECHO(fdSet[i])
			}
		}
	}
}

//func echo(fd int) {
//	defer syscall.Close(fd)
//	var buf [32 * 1024]byte
//	for {
//		nbytes, e := syscall.Read(fd, buf[:])
//		if nbytes > 0 {
//			fmt.Printf(">>> %s", buf)
//			syscall.Write(fd, buf[:nbytes])
//			fmt.Printf("<<< %s", buf)
//		}
//		if e != nil {
//			break
//		}
//	}
//}

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




