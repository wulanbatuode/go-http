package main

// import (
// 	"bufio"
// 	"log"
// 	"net"
// )

// type HttpHead struct {
// 	Method  string
// 	URI     string
// 	Version string
// 	Attr    map[string]interface{}
// }

// type Request struct {
// 	Head HttpHead
// 	Body *[]byte
// }

// func main() {
// 	fd, err := net.Listen("tcp", ":8080")
// 	defer fd.Close()
// 	if err != nil {
// 		log.Fatalln("listen", err)
// 	}
// 	for {
// 		conn, err := fd.Accept()
// 		if err != nil {
// 			log.Println("accept", err)
// 			conn.Close()
// 			continue
// 		}
// 		reader := bufio.NewReader(conn)
// 		headBytes := make([][]byte, 1)
// 		prevPrefix := false
// 		for {
// 			line, prefix, err := reader.ReadLine()
// 			if err != nil {
// 				// conn.Close()
// 				log.Println("readline", err)
// 				break
// 			}
// 			if prevPrefix {
// 				last := &headBytes[len(headBytes)-1]
// 				headBytes[len(headBytes)-1] = append(headBytes[len(headBytes)-1], line)
// 			}
// 			if prefix {
// 			}

// 		}

// 	}

// }
