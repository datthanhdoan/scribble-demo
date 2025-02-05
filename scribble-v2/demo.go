package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	scribble "github.com/sdomino/scribble/v2"
)

// Post kế thừa BaseRecord để có các trường cơ bản
type Post struct {
	scribble.BaseRecord
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	// 1. Khởi tạo Scribble
	db, err := scribble.New("data", nil)
	if err != nil {
		log.Fatalf("Lỗi khi tạo db: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== Demo Scribble v2 ===")
	fmt.Println("Nhấn Enter để tiếp tục mỗi bước...")

	// 2. Tạo một Post mới với ID
	waitForEnter(reader, "\n[Bước 1] Tạo Post mới")

	p := Post{
		BaseRecord: scribble.BaseRecord{
			ID:        "post-1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Title: "Bài viết đầu tiên",
		Body:  "Nội dung bài viết",
	}

	err = db.Write("posts", p)
	if err != nil {
		log.Fatalf("Lỗi khi ghi: %v", err)
	}
	fmt.Printf("✓ Đã tạo Post mới:\n")
	fmt.Printf("  - ID: %s\n  - Title: %s\n  - Body: %s\n  - Created: %v\n\n",
		p.ID, p.Title, p.Body, p.CreatedAt.Format("15:04:05"))

	// 3. Đọc lại Post vừa ghi
	waitForEnter(reader, "[Bước 2] Đọc Post vừa tạo")

	var post Post
	err = db.Read("posts", "post-1", &post)
	if err != nil {
		log.Fatalf("Lỗi khi đọc: %v", err)
	}
	fmt.Printf("✓ Đã đọc Post:\n")
	fmt.Printf("  - ID: %s\n  - Title: %s\n  - Body: %s\n  - Created: %v\n\n",
		post.ID, post.Title, post.Body, post.CreatedAt.Format("15:04:05"))

	// 4. Cập nhật Post
	waitForEnter(reader, "[Bước 3] Cập nhật nội dung Post")

	post.Body = "Nội dung đã được cập nhật"
	post.UpdatedAt = time.Now()
	err = db.Update("posts", post)
	if err != nil {
		log.Fatalf("Lỗi khi cập nhật: %v", err)
	}
	fmt.Printf("✓ Đã cập nhật Post:\n")
	fmt.Printf("  - ID: %s\n  - Title: %s\n  - Body: %s\n  - Updated: %v\n\n",
		post.ID, post.Title, post.Body, post.UpdatedAt.Format("15:04:05"))

	// 5. Đọc lại sau khi cập nhật
	waitForEnter(reader, "[Bước 4] Đọc lại Post sau khi cập nhật")

	var updatedPost Post
	err = db.Read("posts", "post-1", &updatedPost)
	if err != nil {
		log.Fatalf("Lỗi khi đọc: %v", err)
	}
	fmt.Printf("✓ Đã đọc Post sau cập nhật:\n")
	fmt.Printf("  - ID: %s\n  - Title: %s\n  - Body: %s\n  - Created: %v\n  - Updated: %v\n\n",
		updatedPost.ID, updatedPost.Title, updatedPost.Body,
		updatedPost.CreatedAt.Format("15:04:05"),
		updatedPost.UpdatedAt.Format("15:04:05"))

	fmt.Println("=== Demo hoàn tất ===")
}

func waitForEnter(reader *bufio.Reader, message string) {
	fmt.Printf("\n%s\nNhấn Enter để tiếp tục...", message)
	reader.ReadString('\n')
}
