
function isEnter(event) {
    return (event.which && event.which == 13) || (event.keyCode && event.keyCode == 13)
}

$(function(){
    var n = 0;
    var host = "ws://{{.host}}:{{.port}}/echo";
    var conn = new WebSocket(host);
    
    conn.onerror = function(event) {}
    conn.onopen = function(event) {}
    conn.onclose = function(event) {}
    
    conn.onmessage = function(event) {
        json = JSON.parse(event.data)
        console.log(json)
        $("#msg .face").empty().append(json.face)
        $("#msg .emotion").empty().append(json.emotion)
    }
    $(window).unload(function() {
        conn.onclose();
    })
    $("#send").keypress(function (event) {
        if (isEnter(event)) {
            json = JSON.stringify({
                face: $("#send").val(),
                emotion: ""
            })
            conn.send(json);
            $("#send").val('');
        }
    });
})