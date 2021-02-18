package main
import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/web_test/trace"
)

type room struct {
    // 他のクライアントに転送するためのメッセージを保持する
    forward chan []byte
    // チャットルームに参加しようとしている
    join chan *client
    // チャットルームから退出しようとしている
    leave chan *client
    // 全てのクライアント
    clients map[*client]bool
    // tracerはチャットルーム上で行われた捜査のログを受け取る
    tracer trace.Tracer
}

func newRoom() *room {
    return &room{
        forward:    make(chan []byte),
        join:       make(chan *client),
        leave:      make(chan *client),
        clients:    make(map[*client]bool),
        tracer:     trace.Off(),
    }
}

func (r *room) run () {
    for {
        select {
        case client := <-r.join:
            //参加
            r.clients[client] = true
            r.tracer.Trace("人 enter")
        case client := <-r.leave:
            // 退室
            delete(r.clients, client)
            close(client.send)
            r.tracer.Trace("人 exit")
        case msg := <-r.forward:
            r.tracer.Trace("メッセージを受信: ", string(msg))
            // 全てのクライアントにメッセージを転送
            for client := range r.clients {
                select {
                case client.send <- msg:
                    // メッセージを送信
                    r.tracer.Trace(" -- 送信された")
                default:
                    // 送信に失敗
                    delete(r.clients, client)
                    close(client.send)
                    r.tracer.Trace(" -- 送信に失敗。クライアントをクリーンアップします")
                }
            }
        }
    }
}

const (
    socketBufferSize = 1024
    messageBufferSize = 256
)
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    socket, err := upgrader.Upgrade(w, req, nil)
    if err != nil {
        log.Fatal("ServeHTTP:", err)
        return
    }
    client := &client{
        socket: socket,
        send:   make(chan []byte, messageBufferSize),
        room:   r,
    }
    r.join <- client
    defer func() { r.leave <- client }()
    go client.write()
    client.read()
}
