// Global Variables
var sessionId;

function sendDelete() {
	var id = document.getElementById('id');
	console.log("/DSLAM?id=" + id.value);
	var url = "/DSLAM?id=" + id.value;
	var xmlhttp= new XMLHttpRequest();
	xmlhttp.open("DELETE", url, true);
	xmlhttp.send();
	xmlhttp.onreadystatechange=function()
	{
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			document.location.href="/option";
		}
	}
}

function connectToDslam(id) {
	var xmlhttp= new XMLHttpRequest();
	var url = "/session?id=" + id;
	xmlhttp.open("GET", url, true);
	xmlhttp.send();
	xmlhttp.onreadystatechange=function()
	{
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var sessionId = JSON.parse(xmlhttp.responseText).sessionId;
			listCard(sessionId);
		}
	}
}

function listCard(sessionId) {
	var xmlhttp= new XMLHttpRequest();
	var url = "/command?sessionId=" + sessionId + "&command=show equipment slot";
	xmlhttp.open("GET", url, true);
	xmlhttp.send();
	xmlhttp.onreadystatechange=function()
	{
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var commandOut = JSON.parse(xmlhttp.responseText).commandOut;
			alert(commandOut);
		}
	}
}


function getDslam(id) {
	xmlhttp= new XMLHttpRequest();
	var url = "/getDslam?id=" + id;
	xmlhttp.open("GET", url, true);
	xmlhttp.send();
	xmlhttp.onreadystatechange=function()
	{
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var dslam = JSON.parse(xmlhttp.responseText).dslam;
			sessionId = JSON.parse(xmlhttp.responseText).sessionId;
			alert(sessionId);
		}
	}
}
