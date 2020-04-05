package util

import (
	"fmt"
	"syscall"
)

//func main() {
//	//var rfdSet [16]int64
//	//fmt.Println("syscall.Select starting")
//	//for i := range rfdSet {
//	//	rfdSet[i] = 0
//	//}
//	//fmt.Println("FD_ZERO..., rfdSet=", rfdSet)
//	//var fd = 1
//	//rfdSet[fd/64] = 1 << (fd % 64)
//	//fmt.Println("FD_SET..., rfdSet=", rfdSet)
//	var i = 0x01
//	fmt.Println(^(i))
//}

func ECHO(fd int) {
	defer syscall.Close(fd)
	fmt.Println("echo fd:", fd)
	var buf [32 * 1024]byte
	for {
		nbytes, err := syscall.Read(fd, buf[:])
		if err != nil {
			fmt.Println("syscall.Read err=", err)
		}
		if nbytes < 0 {
			if err != syscall.EWOULDBLOCK {
				fmt.Println("nbytes < 0 syscall.Read err=", err)
			}
			//break;
		}
		if nbytes == 0 {
			fmt.Println("Connection closed")
			break
		}

		wn, err := syscall.Write(fd, buf[:nbytes])
		if wn < 0 {
			fmt.Println("syscall.Write err=", err)
			break
		}
		fmt.Printf("gourouting<<< %s \n", buf)


	}
}



/**
假设fd=1是最小的
若fd=1，则fdset：0x0010，若fd=3，则fdset：0x1000，即fd在第几位上
 */
func FD_SET(fd int, p *syscall.FdSet) {
	fmt.Println("before FD_SET, fd=", fd, " rfset=", p)
	p.Bits[fd/64] |= 1 << (uint(fd) % 64)
	fmt.Println("after FD_SET, fd=", fd, " rfset=", p)
}

func FD_ISSET(fd int, p *syscall.FdSet) bool {
	return (p.Bits[fd/64] & (1 << uint(fd) % 64)) != 0
}

func FD_ZERO(p *syscall.FdSet) {
	fmt.Println("before FD_ZERO, p=", p)
	for i := range p.Bits {
		p.Bits[i] = 0
	}
	fmt.Println("after FD_ZERO, p=", p)
}

func FD_CLR(fd int, p *syscall.FdSet) {
	fmt.Println("before FD_CLR, fd=", fd, "    p=", p)
	p.Bits[fd/64] &= (^(1 <<(uint(fd) % 64)))
	fmt.Println("after FD_CLR, fd=", fd, "    p=", p)
}

