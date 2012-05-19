var ws; // WebSocket connection
var conn; // TCP connection
var host;
var port;
var service; // host:port
var status; // "DISCONNECTED", "CONNECTING", "CONNECTED"

$(function () {
	// alert("hello");
	
	$("#connect-button").click(function () {
		status = "DISCONNECTED";
		host = $("#host").val();
		port = $("#port").val();
		if ((port == null) || (port == "")) {
			port = "80";
		}
		if (ws != null) {
			ws.close();
		}
		service = host + ":" + port;
		ws = new WebSocket("ws://localhost:8000/websocket/ws");
		if (ws == null) {
			console.log("WebSocket creation failed");
			return;
		} else {
			console.log("WebSocket creation succeeded");
		}
		$("#console").text("");
		ws.onopen = function (event) {
			console.log("onopen: service=" + service);
			SetStatus("Connecting");
			status = "CONNECTING";
			ws.send(service);
		}
		ws.onerror = function (event) {
			console.log("onerror");
			status = "DISCONNECTED";
			ws.close();
			ws = null;
		}
		ws.onmessage = function (event) {
			if (status == "CONNECTED") {
				console.log("onmessage: " + event.data);
				$("#console").append("<div>" + htmlEncode(event.data) + "</div>");
			} else {
				if (event.data == "SUCC") {
					status = "CONNECTED";
					SetStatus("Connected");
				} else {
					status = "DISCONNECTED";
					SetStatus("Disconnected: " + event.data);
					ws.close();
					ws = null;
				}
			}
		}
		ws.onclose = function (event) {
			console.log("onclose");
			SetStatus("Disconnected");
			status = "DISCONNECTED";
			ws.close();
			ws = null;
		}
	});
	
	function SetStatus(str) {
		$("#status").text(str);
	}
	
	function htmlEncode(value){
		return $('<div/>').text(value).html();
	}

	function SendMessage(message) {
		console.log("Sending message: " + message);
		ws.send(message);
		$("#command").val("");
		
		var dom = $("<div>");
		dom.find("div").append(message);
		console.log("SendMessage: dom is " + dom.html());
		
		$("#console").append("<div class=\"output\">" + htmlEncode(message) + "</div>");
	}
	
	$("#send-button").click(function () {
		var message = $("#command").val();
		SendMessage(message);
	});
	
	$("#command").keypress(function (event) {
		if (event.which == 13) { // Enter key
			var message = $("#command").val();
			SendMessage(message);
		}
	});
});

