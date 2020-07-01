package main

import (
	"fmt"
	"net/rpc"	// provides access to the exported methods of an object across a network or other I/O connection
	 "bufio"	//implements buffered I/O
	 "os"		//provides a platform-independent interface to operating system functionality
	 "strings"	// implements simple functions to manipulate UTF-8 encoded strings
	 "regexp"	//implements regular expression search
	 "strconv"	//implements conversions to and from string representations of basic data types
	 "time"		// to get epoch time
)

type MessageArg struct {
	Chatroom, Chatmsg, Username, Days string
}

type Messages struct {
	ClientMsg []string
}

const test=10

func Split(s string, delim string) []string {
    reg := regexp.MustCompile(delim)
    matches := reg.FindAllStringIndex(s, -1)  //all substrings
    beg := 0
    strings := make([]string, len(matches) + 1)	
    for i, match := range matches {
            strings[i] = s[beg:match[0]]
            beg = match[1]
    }
    strings[len(matches)] = s[beg:len(s)]
    return strings
}

//this will run concurrently with other function
func goRoutineReadMessage(clientRpc *rpc.Client,roomNumber *string,username string,lastMessageId int64, days string){
	for {
		if len(username)>0 {

			//Synchonous Remote call to server method GetChatHistory with argument roomNumber and reply on cliMessages struct type Messages
			cliMessages:=&Messages{}
			err:= clientRpc.Call("Message.GetChatHistory", roomNumber, &cliMessages)	
			if err != nil {
				fmt.Println("error:", err)
			} else {
				receivedMessages:=cliMessages.ClientMsg
				for v := range receivedMessages {
					splitMessages:=Split(receivedMessages[v],"-")

  					//ParseInt interprets a string splitMessages[0] in base 10 and returns corresponding value with bitsize int64. 
					messageIdNumber, _ :=strconv.ParseInt(splitMessages[0], 10, 64)
					
					now := time.Now()
					curSecs := now.Unix()
					days := curSecs / (60 * 60 * 24)
					lastdays, _ :=strconv.ParseInt(splitMessages[4], 10, 64)
					//fmt.Println("Current Day:", days, " lastdays:", lastdays)	

					var datediff int64
					datediff= days - lastdays
					//datediff=test


					if(*roomNumber==splitMessages[1] && lastMessageId<messageIdNumber && username!=splitMessages[3]) { 
						
						lastMessageId = messageIdNumber

						if(datediff < 8){
							fmt.Println("")	
							user := strings.TrimSpace(splitMessages[3])
							fmt.Println(user+": "+splitMessages[2])	
							
						}else{
							splitMessages[2]=""
							fmt.Println("Messages deleted being Older than 7 days!")
							
						}
						
					}
				}
			}
		}	
	}
}



func main() {

	//create inputreader as a pointer to a Reader in bufio. 
	//bufio.NewReader() constructor takes as argument any object that satisfies the io.Reader interface 
	//and returs a new buffered io.Reader that reads from the given reader, os.Stdin satisfies this requirement.
	inputreader := bufio.NewReader(os.Stdin)
	
	logged:=false
	var room="1"
	userInput:=""
	var lastMessageId int64
	lastMessageId=0
	var username string
	var days="0"

	// clients can see the service, to invoke that it first dials the server 
	//only then it will make a remote call
	clientRpc, err := rpc.Dial("tcp", "127.0.0.1:4321")
	if err != nil {
		fmt.Println("dialing:", err)
	}

	fmt.Println("",userInput)
	msgArgs:=&MessageArg{room,"","",""}

	for {
		//check if the user is new by checking logged flag
		if logged==false {
			fmt.Print("Your name:")	
			userInput, _ := inputreader.ReadString('\n')
			userInput=strings.TrimSpace(userInput)
			username=userInput

			//set the flag for user to be true and go to the goRoutineReadMessage method which runs concurrently and will remotely call server method GetChatHistory
			logged=true		
			go goRoutineReadMessage(clientRpc,&room,username,lastMessageId,days)
		} 	else {//ReadString reads until the first occurrence of \n, returning a string containing the data up to and including the newline
				userInput, err := inputreader.ReadString('\n')
				userInput=strings.TrimSpace(userInput)
				startsWith := strings.HasPrefix(userInput, "switch")

			if startsWith==true {
				//replace all "switch" with "" at userInput by puting n=-1
				// If n < 0, there is no limit on the number of replacements.
				room= strings.Replace(userInput, "switch", "", -1)
				fmt.Println("Switching room...\n")
				fmt.Println("You just switched to room :", room)

			} else {
				//Synchonous Remote call to server method OnChat with arguments and reply both as type MessageArg struct
				msgArgs = &MessageArg{room,userInput,username,days} 
				err = clientRpc.Call("Message.OnChat", &msgArgs ,&msgArgs)
				if err != nil {
					fmt.Println("error:", err)
				} else {
					//Call method waits for the remote call to complete
				}
			}
			
		}
	}
}