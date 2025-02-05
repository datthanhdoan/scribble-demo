package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"

    scribble "github.com/sdomino/scribble"
)

// Post là struct đại diện cho dữ liệu (tiêu đề + nội dung)
type Post struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}

// db sẽ là biến toàn cục, lưu kết nối đến Scribble
var db *scribble.Driver

// initDB khởi tạo Scribble
func initDB() {
    var err error
    // Thư mục "data" là nơi Scribble lưu file JSON
    db, err = scribble.New("data", nil)
    if err != nil {
        log.Fatal("Không thể tạo Scribble:", err)
    }
    log.Println("Đã khởi tạo Scribble, dữ liệu lưu trong thư mục 'data/'")
}

// homeHandler xử lý việc hiển thị danh sách Post
func homeHandler(w http.ResponseWriter, r *http.Request) {
    // Đọc tất cả bản ghi trong collection "posts"
    records, err := db.ReadAll("posts")
    if err != nil {
        http.Error(w, "Không thể đọc dữ liệu", http.StatusInternalServerError)
        return
    }

    // Duyệt qua records (dạng []string là JSON thô),
    // parse từng JSON vào struct Post
    var posts []Post
    for _, rec := range records {
        var p Post
        if err := json.Unmarshal([]byte(rec), &p); err == nil {
            posts = append(posts, p)
        }
    }

    // Dữ liệu gửi sang template
    data := struct {
        Posts []Post
    }{
        Posts: posts,
    }

    // Render trang HTML
    t, _ := template.New("home").Parse(homeTemplate)
    t.Execute(w, data)
}

// createHandler hiển thị form tạo Post mới
func createHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.New("create").Parse(createTemplate)
    t.Execute(w, nil)
}

// createPostHandler xử lý khi người dùng submit form
func createPostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Lấy dữ liệu từ form
    title := r.FormValue("title")
    body := r.FormValue("body")

    if title == "" {
        http.Error(w, "Title không được để trống", http.StatusBadRequest)
        return
    }

    // Tạo Post và ghi vào Scribble
    newPost := Post{
        Title: title,
        Body:  body,
    }

    // Ghi vào collection "posts" với key là title
    if err := db.Write("posts", title, newPost); err != nil {
        http.Error(w, "Không thể ghi dữ liệu", http.StatusInternalServerError)
        return
    }

    // Sau khi ghi xong, quay về trang chủ
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HTML template cho trang chủ
var homeTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Danh sách Post</title>
</head>
<body>
    <h1>Danh sách Post</h1>
    <a href="/create">Tạo Post mới</a>
    <ul>
        {{range .Posts}}
        <li><strong>{{.Title}}</strong>: {{.Body}}</li>
        {{end}}
    </ul>
</body>
</html>
`

// HTML template cho trang tạo Post
var createTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Tạo Post mới</title>
</head>
<body>
    <h1>Tạo Post mới</h1>
    <form action="/create-post" method="POST">
        <label>Title: <input type="text" name="title"></label><br><br>
        <label>Body: <textarea name="body" rows="4" cols="40"></textarea></label><br><br>
        <button type="submit">Tạo</button>
    </form>
</body>
</html>
`

func main() {
    // Khởi tạo Scribble
    initDB()

    // Thiết lập router
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/create", createHandler)
    http.HandleFunc("/create-post", createPostHandler)

    fmt.Println("Server đang chạy tại http://localhost:8080 ...")
    // Lắng nghe và phục vụ trên cổng 8080
    log.Fatal(http.ListenAndServe(":8080", nil))
}
