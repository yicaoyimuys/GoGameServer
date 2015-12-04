package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.ListenAndServe(":7224", nil)
}

func wsHandler(ws *websocket.Conn) {
	for {
		s := ""
		if e := websocket.Message.Receive(ws, &s); e != nil {
			break
		}

		fmt.Println(s)

		websocket.Message.Send(ws, []byte("yangsong"))
	}
	ws.Close()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	homeTpl.Execute(w, r.Host)
}

var homeTpl = template.Must(template.New("ws").Parse(`<html>
<textarea id=idout rows=24 cols=72></textarea><hr>
<input id=idin type=search placeholder='Enter a DOS command '
onchange='send(this.value)'></input>
<script>
        var vout=document.getElementById('idout')
        var vin =  document.getElementById('idin')
        var wscon = new WebSocket("ws://{{.}}/ws")
        wscon.onclose = function(e) {vout.value = 'websocket closed'}
        wscon.onmessage = function(e) { vout.value += e.data }
        function send(s) {
                vout.value = ""
                wscon.send(s)
        }
</script>
`))
