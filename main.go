package main

import (
	"bufio"
	"fmt"
	"net"
)

type Room struct {
	users []net.Conn
}

func (r *Room) addUser(user net.Conn) bool {
	if len(r.users) < 5 {
		r.users = append(r.users, user)
		return true
	}
	return false
}

func (r *Room) removeUser(user net.Conn) {
	for i, u := range r.users {
		if u == user {
			r.users = append(r.users[:i], r.users[i+1:]...)
			break
		}
	}
}

func (r *Room) broadcast(msg string, sender net.Conn) {
	for _, user := range r.users {
		if user != sender {
			fmt.Fprintf(user, "[%v] : %s\n", sender.RemoteAddr(), msg)
		}
	}
}

func handleConnection(conn net.Conn, room *Room) {
	defer conn.Close()

	if !room.addUser(conn) {
		fmt.Fprintln(conn, "Room is Full.")
		return
	}

	fmt.Fprintf(conn, "Welcome to the chat room!\n")
	room.broadcast(fmt.Sprintf("user %v has joined the room", conn.RemoteAddr()), conn)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		room.broadcast(text, conn)
		// room.broadcast(fmt.Sprintf("[%v] : %s", conn.RemoteAddr(), text), conn)
	}

	room.removeUser(conn)
	room.broadcast(fmt.Sprintf("[user : %v] has left the chat room.", conn.RemoteAddr()), conn)
}

func main() {
	room := &Room{}

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Chat room started on port 8080.")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, room)

	}

}
