package core

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

//returns a pointer to a socketWorker instance. The instance will already be active on the socket by the time you get it
func GetSocketInstance(callback func(string)) (*socketWorker, error) {
	if callback == nil {
		return nil, errors.New("attempted to create socketWorker without providing read callback")
	}
	sock := &socketWorker{
		running:      true,
		writeChannel: make(chan string, 100),
		readChannel:  make(chan string, 100),
		readCallBack: callback}
	sock.init()
	return sock, nil
}

type socketWorker struct {
	conn         net.Conn
	reader       bufio.Reader
	running      bool
	writeChannel chan string
	readChannel  chan string
	readCallBack func(string)
}

//basically only turns off the socket worker, no current way to restart
func (sock *socketWorker) ToggleRunning() {
	sock.running = !sock.running
}

//initializes the connection, starts the io goroutines
func (sock *socketWorker) init() {
	fmt.Println("Initializing socketWorker...")
	connection, err := net.Dial("tcp", ":2112")
	//@login
	connection.Write([]byte(Login.String()))
	if err != nil {
		fmt.Printf("[socketWorker]\tERROR:%s\n", err)
	}
	//establish connection and creator readers/writers
	sock.conn = connection
	sock.reader = *bufio.NewReader(sock.conn)
	//start working threads
	go sock.listen()
	go sock.run()
}

//handles reading and writing
func (sock *socketWorker) run() {
	defer fmt.Println("[socketWorker]\tWorker thread shut down...")
	fmt.Println("Reading and writing from socket...")
	for sock.running {
		select {
		case message := <-sock.readChannel:
			sock.readCallBack(message)
			break
		case outGoing := <-sock.writeChannel:
			packet := []byte(outGoing)
			if packet[len(packet)-1] != byte('\n') {
				packet = append(packet, byte('\n'))
			}
			sock.conn.Write(packet)
			//fmt.Printf("[socketWorker]\tWriting: %s\n", string(outGoing))
			break
		}
	}
}

//constantly listens for new strings off of the socket and passes them to the read channel
func (sock *socketWorker) listen() {
	fmt.Println("Listening to socket...")
	defer fmt.Println("[socketWorker]\tListening thread shut down...")
	for sock.running {
		message, err := sock.reader.ReadString(byte('\n'))
		if err != nil {
			fmt.Printf("[socketWorker]\tREAD ERROR:%s\n", err)
		}
		sock.readChannel <- message
	}
}

//socketWorker function that will pass strings along to the write channel
func (sock *socketWorker) Write(message string) {
	sock.writeChannel <- message
}
