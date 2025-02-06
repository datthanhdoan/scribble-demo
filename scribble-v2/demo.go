package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	scribble "github.com/sdomino/scribble/v2"
)

// Post là struct cho bài viết
type Post struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	db, reader := setup()
	defer fmt.Println("\n=== Demo kết thúc ===")

	waitForEnter(reader, "\n[Bước 1] Tạo bài viết mới")
	id := createPost(db)

	waitForEnter(reader, "[Bước 2] Đọc bài viết vừa tạo")
	readPost(db, id)

	waitForEnter(reader, "[Bước 3] Cập nhật nội dung bài viết")
	updatePost(db, id)

	waitForEnter(reader, "[Bước 4] Đọc bài viết sau khi cập nhật")
	readPost(db, id)

	waitForEnter(reader, "[Bước 5] Xóa bài viết")
	deletePost(db, id)
}

// Khởi tạo database
func setup() (*scribble.Driver, *bufio.Reader) {
	db, err := scribble.New("data")
	if err != nil {
		log.Fatalf("Lỗi tạo db: %v", err)
	}
	fmt.Println("\n=== Demo Scribble v2 ===")
	fmt.Println("Nhấn Enter để tiếp tục từng bước...")
	return db, bufio.NewReader(os.Stdin)
}

// Tạo bài viết mới
func createPost(db *scribble.Driver) string {
	post := Post{
		Title: "Bài viết đầu tiên",
		Body:  "Nội dung bài viết đầu tiên",
	}

	if err := db.Write("posts", &post); err != nil {
		log.Fatalf("Lỗi ghi: %v", err)
	}

	fmt.Println("✓ Đã tạo bài viết mới")
	return "1" // ID đầu tiên luôn là 1
}

// Đọc bài viết từ database
func readPost(db *scribble.Driver, id string) {
	var post Post
	if err := db.Read("posts", id, &post); err != nil {
		log.Fatalf("Lỗi đọc: %v", err)
	}

	fmt.Printf("\n✓ Đã đọc bài viết:\n")
	fmt.Printf("  - Tiêu đề: %s\n", post.Title)
	fmt.Printf("  - Nội dung: %s\n", post.Body)
}

// Cập nhật bài viết
func updatePost(db *scribble.Driver, id string) {
	post := Post{
		Title: "Bài viết đầu tiên",
		Body:  "Nội dung đã được cập nhật",
	}

	if err := db.Update("posts", id, &post); err != nil {
		log.Fatalf("Lỗi cập nhật: %v", err)
	}

	fmt.Println("✓ Đã cập nhật bài viết")
}

// Xóa bài viết
func deletePost(db *scribble.Driver, id string) {
	if err := db.Delete("posts", id); err != nil {
		log.Fatalf("Lỗi xóa: %v", err)
	}

	fmt.Println("✓ Đã xóa bài viết")
}

func waitForEnter(reader *bufio.Reader, message string) {
	fmt.Printf("\n%s\nNhấn Enter để tiếp tục...", message)
	reader.ReadString('\n')
}
