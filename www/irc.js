/*
 * Wig (web-irc-go)
 * 2016-02-17
*/

//define protobuf & bytebuffer
var ProtoBuf = dcodeIO.ProtoBuf
var ByteBuffer = dcodeIO.ByteBuffer

function IRC() { }
IRC.prototype.Proto = {}

IRC.prototype.OnOpen = function(evt) {
	this.Connect()
}

IRC.prototype.OnMessage = function(evt) {
	var msg = this.Proto.Command.decode(evt.data)
	console.log(msg)
	
	switch(msg.id){
		case 2: {
			
			break
		}
		
	}
}

IRC.prototype.OnError = function(evt) {
	
}

IRC.prototype.OnClose = function(evt) {
	
}

IRC.prototype.Connect = function() {
	var server = "0x.tf"
	var port = 6667
	var ssl = true
	var nick = "testwig"
	
	var msg = new this.Proto.Command({id: 1, connectCommand: {sessionid: "", server: server, port: port, ssl: ssl, nick: nick, realname: nick}})
	this.ws.send(msg.toArrayBuffer())
}

IRC.prototype.Init = function(){
	var self = this //I hate JS
	this.builder = ProtoBuf.loadProtoFile("msgs.proto")
	this.Proto.Command = this.builder.build("main.Command")
	
	this.ws = new WebSocket("ws://" + (window.location.host != "" ? window.location.host : "localhost") + ":9002/ws", "irc")
	this.ws.binaryType = 'arraybuffer';
	this.ws.onopen = function(evt) { self.OnOpen(evt) }
	this.ws.onmessage = function(evt) { self.OnMessage(evt) }
	this.ws.onerror = function(evt) { self.OnError(evt) }
	this.ws.onclose = function(evt) { self.OnClose(evt) }
	
	this.connected = false
}