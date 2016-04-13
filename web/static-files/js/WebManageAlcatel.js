// Global Variables
var data;
var ws = new WebSocket("ws://localhost:8080/ws");
ws.onopen = function()
	{
		ws.send("Connection init");
	};

window.onload = function() {
	ws.send("Send message with JS");
	ws.onmessage = function (evt) 
		{ 
			var received_msg = evt.data;
			alert(received_msg);
			ws.send("Send when receive message");
		};
	if (window.location.pathname == "/index" || window.location.pathname == "/") {
		// Update JSON Data File on Server
		xmlhttp2 = new XMLHttpRequest();
		var url2 = "/SITEAPI/update";
		xmlhttp2.open("GET", url2, true);
		xmlhttp2.send();
		document.getElementById('loader').style.display = "inline-block";
		xmlhttp2.onreadystatechange=function()
		{
			if (xmlhttp2.readyState==4 && xmlhttp2.status==200) {
				data = JSON.parse(xmlhttp2.responseText);
				console.log(data);
				displayDslam();
				document.getElementById('loader').style.display = "none";
			}
		};
	}
	document.getElementById('servicesForm').onsubmit = function() {
		document.getElementById('loader').style.display = "inline-block";
		var name = document.getElementById('name').value;
		var internetSwitch = document.getElementById('internetSwitch').checked;
		var voipSwitch = document.getElementById('voipSwitch').checked;
		var iptvSwitch = document.getElementById('iptvSwitch').checked;
		var dslamID = document.getElementById('dslamID').value;
		var slot = document.getElementById('slot').value;
		var portIndex = document.getElementById('portIndex').value;
		var xmlhttp3 = new XMLHttpRequest();
		var url3 = "/SITEAPI/services?portName=" + name + "&internetSwitch=" + internetSwitch + "&voipSwitch=" + voipSwitch + "&iptvSwitch=" + iptvSwitch + "&dslamID=" + dslamID + "&slot=" + slot + "&portIndex=" + portIndex;
		xmlhttp3.open("GET", url3, true);
		xmlhttp3.send();
		xmlhttp3.onreadystatechange=function()
		{
			if (xmlhttp3.readyState==4 && xmlhttp3.status==200)
			{
				data = JSON.parse(xmlhttp3.responseText);
				console.log(data);
				displayCards(dslamID);
				displayPorts(dslamID, slot);
				displayService(dslamID, slot, portIndex);
				document.getElementById('loader').style.display = "none";
			}
		};
    return false;
	};
};

function displayDslam() {
	// Get DSLAM array in data JSON
	var dslam = data.DSLAM;
	// Remove all DSLAM in DSLAM list div
	while (document.getElementById("DSLAMListDiv").hasChildNodes()) {
		document.getElementById("DSLAMListDiv").removeChild(document.getElementById("DSLAMListDiv").lastChild);
	}
	// Add DSLAM button in DSLAM list div
	for (var key in dslam) {
		var button = document.createElement("BUTTON");
		button.className = "list-button pure-button pure-u-1";
		button.onclick = (function() {
    	var dslamId = dslam[key].Id;
    	return function() {
				displayCards(dslamId);
    	};
  	})()
		var span1 = document.createElement("SPAN");
		span1.className = "fa-stack fa-lg";
		var span2 = document.createElement("SPAN");
		span2.className = "textButtonDSLAM";
		span2.appendChild(document.createTextNode(dslam[key].Name));
		var i1 = document.createElement("I");
		i1.className = "fa fa-cube fa-stack-2x";
		span1.appendChild(i1);
		if (dslam[key].Status != "OK") {
				var i2 = document.createElement("I");
				i2.className = "fa fa-close fa-stack-2x";
				span1.appendChild(i2);
				button.disabled = true;
		}
			button.appendChild(span1);
			button.appendChild(span2);
			document.getElementById("DSLAMListDiv").appendChild(button);
	}
}

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
	};
}

function sendCommand(sessionId, command) {
	var xmlhttp= new XMLHttpRequest();
	var url = "/SITEAPI/command?sessionId=" + sessionId + "&command=" + command;
	xmlhttp.open("GET", url, true);
	xmlhttp.send();
	xmlhttp.onreadystatechange=function()
	{
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var commandOut = JSON.parse(xmlhttp.responseText).commandOut;
			return commandOut;
		}
	};
}

function getDslamById(id) {
	// Get DSLAM array in data JSON
	var dslam = data.DSLAM;
	for (var key in dslam){
		if (dslam[key].Id == id) {
			return dslam[key];
		}
	}
}

function getCardBySlot(dslam, slot) {
	// Get DSLAM array in data JSON
	var card = dslam.Card;
	for (var key in card){
		if (dslam.Card[key].Slot == slot) {
			return dslam.Card[key];
		}
	}
}

function getPortByIndex(card, portIndex) {
	// Get DSLAM array in data JSON
	var port = card.Port;
	for (var key in port){
		if (card.Port[key].Index == portIndex) {
			return card.Port[key];
		}
	}
}

function displayCards(dslamID) {
	var dslam = getDslamById(dslamID);
	while (document.getElementById("CardListDiv").hasChildNodes()) {
		document.getElementById("CardListDiv").removeChild(document.getElementById("CardListDiv").lastChild);
	}
	for (var key in dslam.Card){
		if ((dslam.Card[key].Availability == "available") && (dslam.Card[key].Slot != "nt-a") && (dslam.Card[key].Slot != "nt-b")) {
			var button = document.createElement("BUTTON");
			button.id = dslam.Card[key].Slot;
			button.className = "list-button pure-button pure-u-1";
			button.onclick = (function() {
	    	var dslamId = dslamID;
				var slot = dslam.Card[key].Slot;
	    	return function() {
					displayPorts(dslamId, slot);
	    	};
	   	})()
			var span1 = document.createElement("SPAN");
			span1.className = "fa-stack fa-lg";
			var span2 = document.createElement("SPAN");
			span2.className = "textButtonCard";
			span2.appendChild(document.createTextNode(dslam.Card[key].Name.toUpperCase()));
			var i1 = document.createElement("I");
			i1.className = "imageCard fa fa-hdd-o fa-stack-2x";
			span1.appendChild(i1);
			var i2 = document.createElement("I");
			if (dslam.Card[key]["Opers-Status"] == "no") {
				i2.className = "imageCard fa fa-close fa-stack-2x";
				span1.appendChild(i2);
				button.disabled = true;
			} else if (dslam.Card[key]["Opers-Status"] == "yes" && dslam.Card[key]["Error-Status"] != "no-error") {
				i2.className = "imageCard fa fa-exclamation fa-stack-2x";
				span1.appendChild(i2);
				button.disabled = true;
			}
			button.appendChild(span1);
			button.appendChild(span2);
			document.getElementById("CardListDiv").appendChild(button);
		}
	}
}

function displayPorts(dslamID, slot) {
	var dslam = getDslamById(dslamID);
	var card = getCardBySlot(dslam, slot);
	while (document.getElementById("PortListDiv").hasChildNodes()) {
		document.getElementById("PortListDiv").removeChild(document.getElementById("PortListDiv").lastChild);
	}
	for (var key in card.Port){
			var button = document.createElement("BUTTON");
			button.id = card.Port[key].Index;
			button.className = "list-button pure-button pure-u-1";
			button.onclick = (function() {
	    	var dslamId = dslamID;
				var slot = card.Slot;
				var portIndex = card.Port[key].Index;
	    	return function() {
					displayService(dslamId, slot, portIndex);
	    	};
	   	})()
			var span1 = document.createElement("SPAN");
			span1.className = "fa-stack fa-lg";
			var span2 = document.createElement("SPAN");
			span2.className = "textButtonPort";
			if (card.Port[key].Name !== "") {
				span2.appendChild(document.createTextNode(card.Port[key].Name));
			} else {
				span2.appendChild(document.createTextNode(card.Port[key].Index));
			}
			var i1 = document.createElement("I");
			i1.className = "imagePort fa fa-caret-square-o-right";
			span1.appendChild(i1);
			var i2 = document.createElement("I");
			if (card.Port[key]["Adm-State"] == "down") {
				i2.className = "imagePort fa fa-close fa-stack-2x";
				span1.appendChild(i2);
				//button.disabled = true
			} else if (card.Port[key]["Adm-State"] == "up" && card.Port[key]["Opr-State/Tx-Rate-Ds"] != "down") {
				i2.className = "imagePort fa fa-exclamation fa-stack-2x";
				span1.appendChild(i2);
			}
			button.appendChild(span1);
			button.appendChild(span2);
			document.getElementById("PortListDiv").appendChild(button);

	}
}

function displayService(dslamID, slot, portIndex) {
	var dslam = getDslamById(dslamID);
	var card = getCardBySlot(dslam, slot);
	var port = getPortByIndex(card, portIndex);

	if (port.Name !== "") {
		document.getElementById('name').value = port.Name;
	} else {
		document.getElementById('name').value = port.Index;
	}

	document.getElementById('internetSwitch').checked = false;
	document.getElementById('voipSwitch').checked = false;
	document.getElementById('iptvSwitch').checked = false;

	for (var key in port.Service){
		if (port.Service[key].Vlan == "10") {
			document.getElementById('internetSwitch').checked = true;
		} else if (port.Service[key].Vlan == "20") {
			document.getElementById('voipSwitch').checked = true;
		} else if (port.Service[key].Vlan == "30") {
			document.getElementById('iptvSwitch').checked = true;
		}
	}

	document.getElementById('servicesForm').removeChild(document.getElementById('dslamID'));
	var hidden1 = document.createElement("INPUT");
	hidden1.setAttribute("type", "hidden");
	hidden1.value = dslamID;
	hidden1.id = "dslamID";
	document.getElementById('servicesForm').appendChild(hidden1);

	document.getElementById('servicesForm').removeChild(document.getElementById('slot'));
	var hidden2 = document.createElement("INPUT");
	hidden2.setAttribute("type", "hidden");
	hidden2.value = slot;
	hidden2.id = "slot";
	document.getElementById('servicesForm').appendChild(hidden2);

	document.getElementById('servicesForm').removeChild(document.getElementById('portIndex'));
	var hidden3 = document.createElement("INPUT");
	hidden3.setAttribute("type", "hidden");
	hidden3.value = portIndex;
	hidden3.id = "portIndex";
	document.getElementById('servicesForm').appendChild(hidden3);

	document.getElementById('servicesForm').style.display = "block";

}
