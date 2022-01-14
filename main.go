package main

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ipAddressFormat(srcFilePath string) {
	//srcFilePath := "ip.txt"
	SrcFile, err := os.Open(srcFilePath)
	if err != nil {
		fmt.Printf("[-] Failed to Open file,err = %d\n", err)
		return
	}
	defer SrcFile.Close()
	DstFile, err := os.Create("format.txt")
	if err != nil {
		fmt.Printf("[-] Failed to Create file,err = %d\n", err)
		return
	}
	defer DstFile.Close()
	scanner := bufio.NewScanner(SrcFile)
	var wg sync.WaitGroup
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		if strings.Contains(scanner.Text(), "/") != false {
			wg.Add(1)
			//fmt.Println(scanner.Text())
			go func(str string) {
				defer wg.Done()
				ip := strings.Split(str, "/")[0]
				ipPieceArray := strings.Split(ip, ".")
				maskNum, err := strconv.Atoi(strings.Split(str, "/")[1])
				if err != nil {
					panic(err)
				}
				if maskNum >= 24 {
					count := int(math.Pow(2, float64(32-maskNum)))
					//fmt.Printf(ip+"\t%d\t%d\n", maskNum, count)
					for i := 0; i < count-1; i++ {
						ip := fmt.Sprintf("%s.%s.%s.%d\r\n", ipPieceArray[0], ipPieceArray[1], ipPieceArray[2], i)
						//fmt.Printf(ip)
						DstFile.WriteString(ip)
					}
				} else if maskNum >= 16 {
					count := int(math.Pow(2, float64(24-maskNum)))
					for i := 0; i < count; i++ {
						for j := 0; j < 256; j++ {
							ip := fmt.Sprintf("%s.%s.%d.%d\r\n", ipPieceArray[0], ipPieceArray[1], i, j)
							//fmt.Printf(ip)
							DstFile.WriteString(ip)
						}
					}
				} else if maskNum >= 8 {
					count := int(math.Pow(2, float64(16-maskNum)))
					for i := 0; i < count; i++ {
						for j := 0; j < 256; j++ {
							for k := 0; k < 256; k++ {
								ip := fmt.Sprintf("%s.%d.%d.%d\r\n", ipPieceArray[0], i, j, k)
								//fmt.Printf(ip)
								DstFile.WriteString(ip)
							}
						}
					}
				} else {
					count := int(math.Pow(2, float64(8-maskNum)))
					for i := 0; i < count; i++ {
						for j := 0; j < 256; j++ {
							for k := 0; k < 256; k++ {
								for l := 0; l < 256; l++ {
									ip := fmt.Sprintf("%d.%d.%d.%d\r\n", i, j, k, l)
									//fmt.Printf(ip)
									_, _ = DstFile.WriteString(ip)
								}
							}
						}
					}
				}
			}(scanner.Text())
		} else {
			DstFile.WriteString(scanner.Text() + "\r\n")
		}
	}
	wg.Wait()
}

func MultiPortScan(ip string, startPort, endPort int) {
	startTime := time.Now()
	var wg sync.WaitGroup
	for i := startPort; i < endPort; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", ip, port)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Printf("[-] %s is Close or Banned by FireWall\n", address)
				return
			}
			conn.Close()
			fmt.Printf("[+] %s is Open\n", address)
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(startTime) / 1e9
	fmt.Printf("\n[*] Scan Finished After %d seconds\n", elapsed)
}

func SinglePortScan(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("[-] %s is Close or Banned by FireWall\n", address)
		return false
	}
	conn.Close()
	fmt.Printf("[+] %s is Open\n", address)
	return true
}

func AliveScan(port string) {
	ipAddressFormat("ip.txt")
	SrcFile, err := os.Open("format.txt")
	if err != nil {
		fmt.Printf("[-] Failed to Open file,err = %d\n", err)
		return
	}
	defer SrcFile.Close()
	DstFile, err := os.Create("result.txt")
	if err != nil {
		fmt.Printf("[-] Failed to Create file,err = %d\n", err)
		return
	}
	defer DstFile.Close()
	scanner := bufio.NewScanner(SrcFile)
	var wg sync.WaitGroup
	var counter int
	var startTime = time.Now()

	for scanner.Scan() {
		counter++
		wg.Add(1)
		//fmt.Println(scanner.Text())
		go func(ip string) {
			defer wg.Done()
			//address := fmt.Sprintf("%s:3389", ip)
			address := fmt.Sprintf("%s:%s", ip,port)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Printf("[-] %s is Close or Banned by FireWall\n", address)
				return
			}
			conn.Close()
			fmt.Printf("[+] %s is Open\n", address)
			DstFile.WriteString(ip + "\r\n")
		}(scanner.Text())
	}
	wg.Wait()
	if err := scanner.Err(); err != nil {
		fmt.Printf("[-] Failed to Read file,err = %d\n", err)
		return
	}
	elapsed := time.Since(startTime) / 1e9
	fmt.Printf("\n[*] Scan Finished After %d seconds\n", elapsed)
}

func main() {
	fmt.Print("[*] Present By Assembly Hack1ng\n\n")
	//MultiPortScan("192.168.101.1", 20, 30)
	//SinglePortScan("192.168.101.1", 25)
	cmd :=os.Args
	//fmt.Println(len(cmd))
	if len(cmd) < 2 {
		fmt.Printf("[*] example: scanner.exe {port}\n")
		return
	}
	//AliveScan("80")
	AliveScan(cmd[1])
}
