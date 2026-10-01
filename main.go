package main

import (
	"bufio"
	"crypto/tls" //
	"fmt"
	"net"
)

var clients []net.Conn

func main() {
	// tcp listener on port 42069
	// listener, err := net.Listen("tcp", ":42069")
	// the above was used before i decided to fucking encrypt this shit
	// because why the hell not

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		panic(err)
	}
	config := &tls.Config{ // just the configuration for the 100% unnecissary encryption
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", ":42069", config)

	if err != nil { // check if there is an error. if the error is not empty, quit
		panic(err)
	}
	defer listener.Close()
	fmt.Println("\n[i] Listening on :42069")
	for {
		// wait for client(s)
		con, err := listener.Accept()
		if err != nil {
			continue
		}

		cname := name(con) //gather the user's name

		fmt.Printf("[i] C/%v: %v\n", cname, con.RemoteAddr())

		// keep track of clients
		clients = append(clients, con)

		// give each client its own goroutine
		// this allows multiple clients
		go hc(con, cname)
	}
}

func hc(con net.Conn, cname string) {
	defer con.Close()
	// when the loop tries to return to quit, we just
	// kick off all of the clients as well, then we return

	// read one peice of tcp bullshit at a time
	reader := bufio.NewScanner(con)

	for reader.Scan() {
		// capture any incoming text from a user
		msg := reader.Text()
		fmt.Printf("[C/%v] : %v\n", cname, msg) // print the fucking message

		// send the message to every fucking client
		for _, client := range clients {
			if client == con {
				continue // skip sending the message to the sender
			}
			fmt.Fprint(client, "Cl/", cname+": ", msg, "\n")
		}

	}
	if err := reader.Err(); err != nil {
		fmt.Println("Read error: ", err)
	}
	for i, client := range clients {
		if client == con {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	fmt.Printf("\n[C/%v] Disconnected", cname) // it alredy returned, so disconect the user
	// we determine if the error was a read error
	// or a disconnect (EOF) with the above 'if'
}

func name(con net.Conn) string { // will request name for now, but then later branch to
	// an authenticator when i feel like it
	fmt.Fprintln(con, "Name? ")
	reader := bufio.NewScanner(con) // read stdin from the user

	if reader.Scan() {
		return reader.Text()
	}

	return "anonymous"
}
