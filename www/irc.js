/*
 * Wig (web-irc-go)
 * 2016-02-17
*/

//define protobuf & bytebuffer
var ProtoBuf = dcodeIO.ProtoBuf
var ByteBuffer = dcodeIO.ByteBuffer

function IRC() { }

IRC.prototype.Proto = {}

IRC.prototype.ParseLine= function(server, line){
	if(line.indexOf("PING") === 0){
		//auto repond with pong
		this.SendMessage(server, line.replace("PING", "PONG"))
	}else{
		var ls = line.split(" ", 3)
		console.log(ls)
		this.OnPrivmsg(server,line.substring(0, line.length-1))
	}
}

IRC.prototype.OnOpen = function(evt) {
	this.Connect()
}

IRC.prototype.OnPrivmsg = function(server, msg){
	var d = document.querySelector("#chat")
	var nc = document.createElement("div")
	nc.innerHTML = msg
	d.appendChild(nc)
}

IRC.prototype.OnMessage = function(evt) {
	var msg = this.Proto.Command.decode(evt.data)
	console.log(msg)
	
	switch(msg.id){
		case 2: {
			this.ParseLine(msg.serverMessage.server, msg.serverMessage.msg)
			//this.OnPrivmsg(msg.serverMessage.server, msg.serverMessage.msg)
			break
		}
		case 3: {
			switch(msg.statusMessage.msgtype){
				case 3: {
					if(msg.statusMessage.statuscode == 1){
						this.Nick(msg.statusMessage.msg, "testwig")
						this.User(msg.statusMessage.msg, "testwig", "testwig")
					}else if (msg.statusMessage.statuscode == 2){
						
					}
					break
				}
			}
			break
		}
	}
}

IRC.prototype.OnError = function(evt) {
	console.log(evt)
}

IRC.prototype.OnClose = function(evt) {
	console.log(evt)
}

IRC.prototype.Connect = function() {
	var server = "0x.tf"
	var port = 6667
	var ssl = true
	
	var msg = new this.Proto.Command({id: 1, connectCommand: {sessionid: "", server: server, port: port, ssl: ssl}})
	this.ws.send(msg.toArrayBuffer())
}

IRC.prototype.Nick = function(server, user){
	var msg = new this.Proto.Command({id: 2, serverMessage: {server: server, msg: "NICK " + user + "\n"}})
	this.ws.send(msg.toArrayBuffer())
}

IRC.prototype.User = function(server, nick, realname){
	var msg = new this.Proto.Command({id: 2, serverMessage: {server: server, msg: "USER " + nick + " 0 * :" + realname + "\n"}})
	this.ws.send(msg.toArrayBuffer())
}

IRC.prototype.Chat = function() {
	var server = "0x.tf"
	var dv = document.querySelector("#ch")

	this.SendMessage(server, dv.value + "\n")
}

IRC.prototype.SendMessage = function(server, msg){
	var msg = new this.Proto.Command({id: 2, serverMessage: {server: server, msg: msg}})
	this.ws.send(msg.toArrayBuffer())
}

IRC.prototype.Init = function(){
	var self = this //I hate JS
	this.builder = ProtoBuf.loadProtoFile("msgs.proto")
	this.Proto.Command = this.builder.build("main.Command")
	
	this.ws = new WebSocket("wss://" + window.location.host + "/ws", "irc")
	this.ws.binaryType = 'arraybuffer';
	this.ws.onopen = function(evt) { self.OnOpen(evt) }
	this.ws.onmessage = function(evt) { self.OnMessage(evt) }
	this.ws.onerror = function(evt) { self.OnError(evt) }
	this.ws.onclose = function(evt) { self.OnClose(evt) }
	
	this.connected = false
}