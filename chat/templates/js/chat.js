$(function(){
    var socket = null;
    var msgBox = $("#chatbox textarea");
    var messages = $("#messages");
    $("#chatbox").submit(function(){
        if(!msgBox.val()) return false;
        if(!socket) {
            alert("websocket接続されていません");
            return false;
        }
        socket.send(JSON.stringify({"Message": msgBox.val()}));
        msgBox.val("");
        return false;
    });
    if(!window["WebSocket"]) {
        alert("websocket非対応のブラウザ");
    }else {
        socket = new WebSocket("ws://{{.Host}}/room");
        socket.onclose = function() {
            alert("接続終了");
        }
        socket.onmessage = function(e) {
            var msg = eval("("+e.data+")");
            messages.append(
                $("<li>").append(
                    $("<img>").css({
                        width:50,
                        verticalAlign:"middle"
                    }).attr("src", msg.AvatarURL),
                    $("<strong>").text(msg.Name + ": "),
                    $("<span>").text(msg.Message)
                )
            );

        }
    }
});

