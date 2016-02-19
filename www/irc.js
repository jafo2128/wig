/*
 * Wig (web-irc-go)
 * 2016-02-17
*/

//define protobuf & bytebuffer
var ProtoBuf = dcodeIO.ProtoBuf
var ByteBuffer = dcodeIO.ByteBuffer

//Line parser translated from https://github.com/fluffle/goirc
function Line(s) {
	this.Tags = {}
	this.Nick = null;
	this.Ident = null;
	this.Host = null;
	this.Src = null;
	this.Cmd = null;
	this.Raw = null;
	this.Args = [];
	this.Time = null;
	
	this.Parse(s)
}

Line.prototype.Parse = function(s){
	this.Raw = s
	this.Time = new Date().getTime()
	
	if (s[0] == '@') {
		var rawTags
		var idx = s.indexOf(" ")
		if (idx >= 0) {
			rawTags = s.substring(1, idx)
			s = s.substring(idx+1)
		} else {
			return null
		}

		tags = rawTags.split(";")
		for(var x = 0; x < tags.length; x++) {
			var tag = tags[x]
			if (tag == "") {
				continue
			}

			var pair = tag.replace(/[\:|\;|\s| |\r|\n]/g, "").split("=", 2)
			if (pair.length < 2) {
				this.Tags[tag] = ""
			} else {
				this.Tags[pair[0]] = pair[1]
			}
		}
	}

	if (s[0] == ':') {
		var idx = s.indexOf(" ")
		if (idx >= 0) {
			this.Src = s.substring(1, idx)
			s = s.substring(idx+1)
		} else {
			return null
		}

		this.Host = this.Src
		var nidx = this.Src.indexOf("!")
		var uidx = this.Src.indexOf("@")
		if (uidx != -1 && nidx != -1) {
			this.Nick = this.Src.substring(0, nidx)
			this.Ident = this.Src.substring(nidx+1, uidx)
			this.Host = this.Src.substring(uidx+1)
		}
	}

	var args = s.split(" :", 2)
	if (args.length > 1) {
		var msg = args[1]
		args = args[0].trim().split(" ")
		args[args.length] = msg
	} else {
		args = args[0].trim().split(" ")
	}
	this.Cmd = args[0].toUpperCase()
	if (args.length > 1) {
		this.Args = args.slice(1)
	}

	if ((this.Cmd === "PRIVMSG"|| this.Cmd === "NOTICE") && this.Args[1].length > 2 && this.Args[1].startsWith("\x001") && this.Args[1].endsWith("\x001")) {
		// WOO, it's a CTCP message
		var t = this.Args[1].replace("\001" ,"").split(" ", 2)
		if (t.length > 1) {
			this.Args[1] = t[1]
		}
		var c = t[0].toUpperCase()
		if (c === "ACTION" && this.Cmd === "PRIVMSG") {
			this.Cmd = c
		} else {
			if (this.Cmd === "PRIVMSG") {
				this.Cmd = "CTCP"
			} else {
				this.Cmd = "CTCPREPLY"
			}
			this.Args = [c].concat(this.Args)
		}
	}
}

function IRC() { }

IRC.prototype.Proto = {}

IRC.prototype.DecodeLine = function(s) {
	return new Line(s)
}

IRC.prototype.ParseLine= function(server, line) {
	if(line.indexOf("PING") === 0){
		//auto repond with pong
		this.SendMessage(server, line.replace("PING", "PONG"))
	}else{
		var ls = this.DecodeLine(line)
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