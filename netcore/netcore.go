package netcore

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

//returns a pointer to a socketWorker instance. The instance will already be active on the socket by the time you get it
func GetSocketInstance() *socketWorker {
	sock := &socketWorker{
		running: true,
		writeChannel: make(chan string),
		readChannel:  make(chan string)}
	sock.init()
	return sock
}

type socketWorker struct {
	conn         net.Conn
	reader       bufio.Reader
	running      bool
	writeChannel chan string
	readChannel  chan string
}

//test function. call it to test things
func (sock *socketWorker) Test() {
	fmt.Println("Waiting for 1 minute")
	time.Sleep(time.Minute * 1)
	sock.running = false
}

func (sock *socketWorker) ToggleRunning() {
	sock.running = !sock.running
}

//initializes the connection
func (sock *socketWorker) init() {
	connection, err := net.Dial("tcp", ":2112")
	if err != nil {
		fmt.Printf("[socketWorker]\tERROR:%s\n", err)
	}
	//establish connection and creater readers/writers
	sock.conn = connection
	sock.reader = *bufio.NewReader(sock.conn)
	//start working threads
	go sock.listen()
	go sock.run()
}

//handles reading and writing
func (sock *socketWorker) run(){
	defer fmt.Println("[socketWorker]\tWorker thread shut down...")
	for sock.running{
		select {
		case message:=<- sock.readChannel :
			//@todo: make a function to handle io reads
			fmt.Printf("%s\n", message)
			case outGoing:= <-sock.writeChannel:
				packet:= []byte(outGoing)
				if packet[len(packet)-1] != byte('\n') {
					packet = append(packet, byte('\n'))
				}
				sock.conn.Write(packet)
				fmt.Printf("[socketWorker]\tWriting: %s", string(outGoing))
		}
	}
}

//constantly listens for new strings off of the socket and passes them to the read channel
func (sock *socketWorker) listen() {
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
