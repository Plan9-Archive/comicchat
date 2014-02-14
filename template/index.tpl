<!doctype html>

<head>
	<title>comicchat</title>
	<script src="https://code.jquery.com/jquery-2.0.0b1.js"></script>
	<script src="//netdna.bootstrapcdn.com/bootstrap/3.1.1/js/bootstrap.min.js"></script>
	<link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap.min.css">
	<style type="text/css">
	</style>

	<script>
	var cc = {};

	cc.connect = function() {
		$("#status").removeClass("btn-danger").removeClass("btn-success").
			addClass("btn-info").html("connecting...");

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
				$("#status").removeClass("btn-info").removeClass("btn-success").
					addClass("btn-danger").html("connection broke");
					window.setTimeout(cc.connect(), 1000);
			}
			cc.conn.onerror = function(evt) {
				$("#status").removeClass("btn-info").removeClass("btn-success").
					addClass("btn-danger").html("error: " + evt.data);
			}
			cc.conn.onmessage = function(evt) {
				msg = JSON.parse(evt.data);
				switch(msg.Type) {
				case "newimage":
					$("#output").append("<img src='" + msg.Message + "' />");

					$('html, body').animate(
					{scrollTop: $('#output').height() + ($("#navbar").height()*2)},
					100
					);

					break;
				case "connected":
					$("#status").removeClass("btn-info");
					$("#status").removeClass("btn-danger");
					$("#status").addClass("btn-success");
					$("#status").html("connection ok!");
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
<div id="output" style="padding-bottom: 50px;"></div>
<!--
<input style="display:none" id="name" /><button style="display:none" id="set">set name</button>
-->
<nav id="navbar" class="navbar navbar-default navbar-left navbar-fixed-bottom" role="navigation">
	<div class="container-fluid">
		<form class="form-inline">
			<div class="form-group">
				<input type="text" class="form-control" id="text">
				<button id="send" class="btn btn-primary">send</button>
			</div>
			<div class="form-group">
				<label class="sr-only" for="face">face:</label>
				<select id="face" class="btn btn-default">
				{{range .Faces}}
					<option value="{{.}}">{{.}}</option>
				{{end}}
				</select>
			</div>

			<div class="form-group">
				<input type="text" class="form-control" id="nick">
				<button id="setnick" class="btn btn-primary">set nick</button>
			</div>
			<div class="form-group pull-right">
				<button id="status" class="btn btn-info navbar-right">connecting...</button>
			</div>
		</form>
	</div>
</nav>
</body>
</html>

