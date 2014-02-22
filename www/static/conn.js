
function isEnter(event) {
    return (event.which && event.which == 13) || (event.keyCode && event.keyCode == 13)
}

$(function(){
    var n = 0;
    var host = "ws://{{.host}}:{{.port}}/echo";
    var conn = new WebSocket(host);
    
    conn.onerror = function(e) {
        alert(e)
    }
    conn.onopen = function() {
        conn.send('in');
    }
    conn.onclose = function() {
        conn.send('quit');
    }
    conn.onmessage = function(event) {
        $("#msg").empty()
        $("#msg").append(event.data + "<br>")
    }
    $(window).unload(function() {
        conn.onclose();
    })
    $("#send").keypress(function (event) {
        if (isEnter(event)) {
            conn.send($("#send").val());
            $("#send").val('')
        }
    });
})