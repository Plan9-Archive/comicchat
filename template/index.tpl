<!doctype html>

<head>
	<title>comicchat</title>
	<script src="https://code.jquery.com/jquery-2.0.0b1.js"></script>
	<script src="//netdna.bootstrapcdn.com/bootstrap/3.1.1/js/bootstrap.min.js"></script>
	<link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap.min.css">
	<style type="text/css">
	body {
		padding-bottom:100px;
	}
	</style>

	<script>
	var cc = {};

	cc.connect = function() {
		// image urls over websocket
		var loc = window.location, new_uri;
		if (loc.protocol === "https:") {
			new_uri = "wss:";
		} else {
			new_uri = "ws:";
		}
		new_uri += "//" + loc.host;
		new_uri += loc.pathname + "new";

		if (window["WebSocket"]) {
			cc.conn = new WebSocket(new_uri);
			cc.conn.onclose = function(evt) {
			}
			cc.conn.onmessage = function(evt) {
				msg = JSON.parse(evt.data);
				switch(msg.Type) {
				case "newimage":
					$("#output").append("<img src='" + msg.Message + "' />");

					$('html, body').animate(
					{scrollTop: $('#output').height()},
					100
					);

					break;
				}
			}
		}
	};

	cc.getName = function () {
		$("#name").show();
		$("#set").show();
		$("#set").click(function () {
			localStorage['comicchat_name'] = $("#name").val();
			$("#name").hide();
			$("#set").hide();
		});
	};

	cc.send = function(typ, msg) {
		var cmsg = {};
		cmsg.Type = typ;
		cmsg.Message = msg;
		cc.conn.send(JSON.stringify(cmsg));
	};


	$(document).ready(function () {
		cc.connect()
		$("#text").keypress(function (e) {
			if (e.which == 13) {
				$("#send").focus().click();
				$("#text").val('');
				$("#text").focus()
				return false;    //<---- Add this line
			}
		});

		$("#send").click(function () {
			cc.send("privmsg", $("#text").val());
			$("#text").val('');
			$("#text").focus()
			return false;    //<---- Add this line
		});

		$("#face").on('change', function() {
			cc.send("action", $("#face").val());
			return false;    //<---- Add this line
		});

		$("#setnick").click(function () {
			cc.send("nick", $("#nick").val());
			$("#text").focus()
			return false;    //<---- Add this line
		});

		cc.name = localStorage['comicchat_name'];
		if (! cc.name) {
			cc.getName();
		}

	});

</script>
</head>
<body>
<div id="output"></div>
<!--
<input style="display:none" id="name" /><button style="display:none" id="set">set name</button>
-->
<nav class="navbar navbar-default navbar-left navbar-fixed-bottom" role="navigation">
	<form class="form-inline">
		<input type="text" class="form-control" width=80 id="text">
		<button id="send" class="btn btn-default">send</button>

		<label for="face">face:</label>
		<select id="face" class="btn btn-default">
		{{range .Faces}}
			<option value="{{.}}">{{.}}</option>
		{{end}}
		</select>

		<input type="text" class="form-control" width=80 id="nick">
		<button id="setnick" class="btn btn-default">set nick</button>
	</form>
</nav>
</body>
</html>

