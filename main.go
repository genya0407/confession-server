package main

import (
	"fmt"
	auth "github.com/genya0407/confession-server/authorization"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func WebSock(w http.ResponseWriter, r *http.Request, ps httprouter.Params, t auth.SessionToken) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(t))
	if err != nil {
		log.Println(err.Error())
		return
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(msg))
}

func HelloInternal(w http.ResponseWriter, r *http.Request, ps httprouter.Params, t auth.SessionToken) {
	fmt.Fprintf(w, "%s, %s!\n", ps.ByName("greet"), t)
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/hellointernal/:greet", auth.AuthorizeBearer(HelloInternal))
	router.GET("/connect", auth.AuthorizeBearer(WebSock))

	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
