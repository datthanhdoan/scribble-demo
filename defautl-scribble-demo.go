package main

import (
    "fmt"
    "log"

    scribble "github.com/sdomino/scribble"
)

// Post là một struct ví dụ để lưu trữ
type Post struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}

func main() {
    // 1. Khởi tạo Scribble, chỉ định thư mục sẽ lưu file JSON
    //    "data" là tên folder sẽ chứa dữ liệu
    db, err := scribble.New("data", nil)
    if (err != nil) {
        log.Fatalf("Error when creating db: %v", err)
    }

    // 2. Tạo một đối tượng Post
    p := Post{"MyTitle", "MyContent"}

    // 3. Ghi đối tượng p vào collection "posts"
    //    với key = p.Title (MyTitle)
    err = db.Write("posts", p.Title, p)
    if (err != nil) {
        log.Fatalf("Error when writing: %v", err)
    }
    fmt.Println("Đã ghi dữ liệu cho Post:", p.Title)

    // 4. Đọc lại dữ liệu vừa ghi
    var post Post
    err = db.Read("posts", "MyTitle", &post)
    if (err != nil) {
        log.Fatalf("Error when reading: %v", err)
    }
    fmt.Printf("Đã đọc dữ liệu: %+v\n", post)
}
