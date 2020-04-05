package main

import (
	"fmt"
	"go-epoll/go-epoll/com/zzz/util"
	"net"
	"os"
	"syscall"
	"time"
)

func main() {

	fd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)


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

	masterSet := syscall.FdSet{}
	util.FD_ZERO(&masterSet)
	util.FD_SET(fd, &masterSet)
	var maxFd = fd
	for {
		var workingSet = syscall.FdSet{Bits: masterSet.Bits}
		fmt.Println("before select", workingSet)
		//待测试的描述集总是从0， 1， 2， …开始，如果maxFd=4，则需要监控5个文件描述符
		timeout := syscall.Timeval{2, 0}
		n, err := syscall.Select(maxFd+1, &workingSet, nil, nil, &timeout)
		fmt.Println("after select ", workingSet)
		time.Sleep(time.Second)
		if err != nil {
			fmt.Println("syscall.Select, err=", err)
			break
		}
		if n < 0 {
			fmt.Println("syscall.Select，, n=", n)
			break
		}
		fdSize := min(syscall.FD_SETSIZE, maxFd)
		for i:=0; i <= fdSize; i++ {
			if util.FD_ISSET(i, &workingSet) {
				if i == fd {
					for {
						nfd, _, err := syscall.Accept(fd)
						fmt.Println("syscall.Accept, fd=", fd, "  nfd=", nfd, "    err=", err)
						if nfd < 0 {
							if err != syscall.EWOULDBLOCK {
								fmt.Println("syscall.Accept err=", err)
								os.Exit(0)
							}
							break
						}
						fmt.Println("a new connect create, nfd=", nfd)
						util.FD_SET(nfd, &masterSet)
						if nfd > maxFd {
							maxFd = nfd
						}
					}
				}else {
					fmt.Println("do else...fd=", i)
					time.Sleep(time.Second*2)
					go func(fd int) {
						fmt.Println("open a gorouting, fd=", fd)
						util.ECHO(fd)
					}(i)
					fmt.Println("do else clr...fd=", i)
					util.FD_CLR(i, &masterSet) //因为使用gorouting来处理，需要把文件描述符移除，否则对一个socket会开多个gorouting
					if i == maxFd {
						for tmp:=i-1; tmp > 0; tmp-- {
							if util.FD_ISSET(tmp, &masterSet) {
								maxFd = tmp
								fmt.Println("change maxFd, maxFd=", maxFd)
							}
						}
					}
					fmt.Println("do else maxfd=", maxFd, " masterSet:", masterSet)
					//也可以不使用gorouting处理，这样就单线程处理所有请求，方式：对每个socket进行读取，也不需要FD_CLD，移除的话，否则下次就监听不到了
					/**
					for {
						var buf [32 * 1024]byte
						nbytes, err := syscall.Read(i, buf[:])
						if nbytes == -1 {	//表示读完
							if err != nil && err != syscall.EWOULDBLOCK {
								util.FD_CLR(i, &masterSet)
								close(i)
							}
							break
						}
						if nbytes == 0 {
								util.FD_CLR(i, &masterSet)
								close(i)
								break;
						}
						wn, err := syscall.Write(i, buf[:nbytes])
					}
					 */
				}
			}
		}
		fmt.Println("after select handle maxFd=", maxFd, "    workingSet=", workingSet)
		//time.Sleep(time.Second * 2)
	}
}

func min(a int, b int) int {
	if a > b {
		return b
	}else{
		return a
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


