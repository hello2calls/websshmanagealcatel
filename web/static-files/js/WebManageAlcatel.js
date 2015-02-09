function sendDelete() {
	var id = document.getElementById('id');
	console.log("/DSLAM?id=" + id.value);
	var url = "/DSLAM?id=" + id.value;
	xmlhttp= new XMLHttpRequest();
	xmlhttp.open("DELETE", url, true);
	xmlhttp.send();
	document.location.href="/option";
}
