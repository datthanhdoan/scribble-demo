package main

import (
    "fmt"
    "log"

    scribble "github.com/sdomino/scribble/v2" // using local upgraded version
)

// Post is an example struct for storage
type Post struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}

func main() {
    // 1. Initialize Scribble, specify directory to store JSON files
    //    "data" is the folder name that will contain the data
    db, err := scribble.New("data", nil)
    if (err != nil) {
        log.Fatalf("Error when creating db: %v", err)
    }

    // 2. Create a Post object
    p := Post{"MyTitle", "MyContent"}

    // 3. Write object p to "posts" collection
    //    with key = p.Title (MyTitle)
    err = db.Write("posts", p.Title, p)
    if (err != nil) {
        log.Fatalf("Error when writing: %v", err)
    }
    fmt.Println("Data written for Post:", p.Title)

    // 4. Read back the written data
    var post Post
    err = db.Read("posts", "MyTitle", &post)
    if (err != nil) {
        log.Fatalf("Error when reading: %v", err)
    }
    fmt.Printf("Data read: %+v\n", post)
}
