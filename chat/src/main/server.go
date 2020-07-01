
package main

import (
	"fmt"
 	"net"
 	"net/rpc"
 	"container/list"    //implements a doubly linked list.
 	"strconv"
 	//"database/sql"		//provides a generic interface around SQL (or SQL-like) databases
 	//_"master/mysql"		//MySQL driver for Go's database/sql package
 	"time"				//for epoch time
 	//"math"
 )

type Message struct {
	Chatroom, Chatmsg, Username, Days string
}

type Messages struct {
	ClientMsg []string
}

var stringChatroomMap map[string]*list.List
var messageId int

//create a exported method "OnChat" for type Client
//the method has two arguments, both built-in/exported types.
//first argument is provided by the caller, which is a struct pointer *Message
//the method's second argument is a struct pointer *Message, is the result parameters to be returned to the caller
//the method has return type error -- (string-->send by errors.New) only if error is returned then the reply parameter will not be sent back to the client.
//server may handle requests on a single connection by calling ServeConn. 
//more typically it will create a network listener and call Accept
//to use this service (object), function Dial on client establishes a connection and then invokes NewClient on the connection
//The resulting Client object has two methods (Call and Go), a pointer containing the arguments, and a pointer to receive the result parameters.
//It hasn't used the Call method which basically waits for the remote call to complete (synchonous)
//It rather uses the Go method that launches the call asynchronously.

func (cl *Message) OnChat(clientmsg *Message, replyclient *Message) error {
	var msg string
	//creating one string for each message by appending time, and it is seperated by - character
	now := time.Now()
	curSecs := now.Unix()
	days := curSecs / (60 * 60 * 24)
	clientmsg.Days=strconv.FormatInt(days,10)
	msg+=strconv.Itoa(messageId)+"-"+clientmsg.Chatroom+"-"+clientmsg.Chatmsg+"-"+clientmsg.Username+"-"+clientmsg.Days
	
/*
	//storing message to mysql database.
	//BEGIN	
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/gochat")
	if err != nil {
		//log.Fatal(err)
		fmt.Println("error")
	}

	dato,err2 := db.Query("insert into chat_log (id, user, text, roomnumber, date) values ('', '"+clientmsg.Username+"', '"+clientmsg.Chatmsg+"','"+clientmsg.Chatroom+"','')")
	//dato,err2 := db.Query("insert into chat_log (id, login, message, room, createdate) values ('', '"+msgs.Login+"', '"+msgs.Texts+"',	'"+msgs.Rooms+"','')")
	if err2!=nil{
		fmt.Println("error2")
	}
	for dato.Next(){
		//var id,login,message,room,createdate string
		err3 := dato.Scan()
		if err3 != nil {
			fmt.Println("error3")
		}
		//fmt.Println("for room "+room+"\n \t id:"+id+" From:"+login+" Message:"+message)
	}
	//END
*/
	//storing message to stringChatroomMap (<chatRoom> <string list>)
	stringChatroomMap[clientmsg.Chatroom].PushFront(msg)
	fmt.Println("Message from user [",clientmsg.Username + " ] for room " + clientmsg.Chatroom + ": " + clientmsg.Chatmsg)
	messageId++
	return nil
}

//create a exported method "GetChatHistory" for type Client
//service (server registred object) is *Message
//caller sends RoomNumber
//returns *Messages as the reuslt parameter to the caller

func (cl *Message) GetChatHistory(RoomNumber string, cliMessages *Messages) error {
	messages := []string{}
	l := stringChatroomMap[RoomNumber]

	/*now := time.Now()
	curSecs := now.Unix()
	days := curSecs / (60 * 60 * 24)

	clientmsg.Days

	datediff = days - 

	messageIdNumber, _ :=strconv.ParseInt(splitMessages[0], 10, 64)*/


	for el := l.Front(); el != nil; el = el.Next() {
	//	fmt.Println(el.Value)
		msgStr:=el.Value.(string)

		messages = append([]string{msgStr}, messages...)
	}
	cliMessages.ClientMsg=messages
	return nil
}


func main() {

	fmt.Printf("Waiting for clients...\n")
	stringChatroomMap = make(map[string]*list.List)
	for i :=1; i<21; i++ {
		stringChatroomMap[strconv.Itoa(i)]= list.New()
	}
		
	messageId=0
	
	messageNew := new(Message)

	rpc.Register(messageNew)

	ln, err := net.Listen("tcp",":4321")
	if err != nil {
		fmt.Printf("Error during server creation ",err)
		return
	}
	for {
		if conn, er := ln.Accept(); er != nil {
			fmt.Printf("Error caught ", er.Error())
		} else {
			now := time.Now()
			curSecs := now.Unix()	   // get elapsed time since the Unix epoch in seconds
			fmt.Println("new connection established on", now.Format("20060102150405"),"since epoch:",curSecs, "seconds!")
			
			go rpc.ServeConn(conn)    //go handleConnection(conn)
		}
	}
}

