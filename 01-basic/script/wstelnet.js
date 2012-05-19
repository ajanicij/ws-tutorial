var host;
var port;
var service; // host:port

$(function () {
	
	$("#connect-button").click(function () {
		alert("You clicked 'Connect'");
		SetStatus("Connecting");
	});
	
	function SetStatus(str) {
		$("#status").text(str);
	}

	function htmlEncode(value){
		return $('<div/>').text(value).html();
	}

	function SendMessage(message) {
		SetStatus("Sending message: " + message);
		console.log("Sending message: " + message);
		$("#command").val("");
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

