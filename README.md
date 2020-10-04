# kurakani app - Chat app using Go Lang

  Package net/rpc has been used to get the portable interface for network I/O (TCP/IP). The console application has used only the basic interface provided by the Dial, Listen, and Accept functions and the associated Conn and Listener interfaces. 
  
  The Listen function creates server at port 4321 and the Dial function from the client connects to the server. 
  
  The message and the chatroom for the specific user is stored in the stringChatroomMap as doubly link list. Mysql database can also be used to keep the storage of the messages. 
  
  Server has two exported procedure OnChat and GetChatHistory which are remotely called by the client.
  
  As the server starts it waits for the connection with clients. Once the client gets connected by dialing at port 4321, it allows user to enter the name and then the new connection will be established with the server. It will then invoke goroutine ReadMessage that runs concurrently and remotely calls server method to get the chat history for the corresponding room.
  
  The user can change the room by typing “switch”, eg: switch6 will take the user to chatroom 6. If the room was active with previous messages, then all of the messages will be displayed to the new user.
  
  The client allows user to input the messages and that will be understood by the server by reading newline character.
  
  When the room remains inactive for 7 days, then the messages stored in the room gets removed using epoch time concept from the time package to calculate the difference in the current time and the last time when the room was active. Since the message is stored in the map, the server should be on all the time, and to test the mechanism of deleting messages after 7 days, const test is assigned an initial value of 10 and used that value at the condition checking block in the program.
