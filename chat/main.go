package main

import(
    "github.com/stretchr/gomniauth"
    "github.com/stretchr/objx"
    "github.com/stretchr/gomniauth/providers/facebook"
    "github.com/stretchr/gomniauth/providers/github"
    "github.com/stretchr/gomniauth/providers/google"
    "os"
    "github.com/web_test/trace"
    "log"
    "net/http"
    "text/template"
    "path/filepath"
    "sync"
    "flag"
)

type templateHandler struct {
    once        sync.Once
    filename    string
    templ       *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t.once.Do(func() {
        t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
    })
    data := map[string]interface{}{
        "Host": r.Host,
    }
    if authCookie, err := r.Cookie("auth"); err == nil {
        data["UserData"] = objx.MustFromBase64(authCookie.Value)
    }

    t.templ.Execute(w, data)
}

func main() {
    const google_ID string = "566614770599-sk13e2v2ru6m6tnbhn91nsmjv021f3ki.apps.googleusercontent.com"
    const google_secret string = "l03SAtcO50kzJgRUhJZdpq_I"

    var addr = flag.String("addr", ":8080", "アプリのアドレス")
    flag.Parse() // フラグを解釈

    // Gomniauthのセットアップ
    gomniauth.SetSecurityKey("セキュリティキー")
    gomniauth.WithProviders(
        facebook.New("aa", "秘密の値", "http://localhost:8080/auth/callback/facebook"),
        github.New("クライアントID", "秘密の値", "http://localhost:8080/auth/callback/github"),
        google.New( google_ID, google_secret, "http://localhost:8080/auth/callback/google"),
    )

    r := newRoom()
    r.tracer = trace.New(os.Stdout)
    http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
    http.Handle("/login", &templateHandler{filename: "login.html"})
    http.Handle("/room", r)
    http.HandleFunc("/auth/", loginHandler)
    http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
        http.SetCookie(w, &http.Cookie{
            Name: "auth",
            Value: "",
            Path: "/",
            MaxAge: -1,
        })
        w.Header()["Location"] = []string{"/chat"}
        w.WriteHeader(http.StatusTemporaryRedirect)
    })
    // チャットルームを開始
    go r.run()

    // webサーバーを起動
    log.Println("webサーバーを開始　ポート: ", *addr)
    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

