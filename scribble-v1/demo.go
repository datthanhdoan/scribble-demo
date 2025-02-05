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
	if err != nil {
		log.Fatalf("Lỗi khi tạo db: %v", err)
	}

	// 2. Tạo một đối tượng Post
	p := Post{
		Title: "Bài viết đầu tiên",
		Body:  "Nội dung bài viết",
	}

	// 3. Ghi đối tượng p vào collection "posts"
	//    với key = p.Title
	err = db.Write("posts", p.Title, p)
	if err != nil {
		log.Fatalf("Lỗi khi ghi: %v", err)
	}
	fmt.Printf("Đã ghi dữ liệu cho Post: %s\n", p.Title)

	// 4. Đọc lại dữ liệu vừa ghi
	var post Post
	err = db.Read("posts", "Bài viết đầu tiên", &post)
	if err != nil {
		log.Fatalf("Lỗi khi đọc: %v", err)
	}
	fmt.Printf("Đã đọc dữ liệu: %+v\n", post)
}
