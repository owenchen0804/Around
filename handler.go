package main
import (
    "encoding/json"
    "fmt"
    "net/http"
	"path/filepath"

    "github.com/pborman/uuid"
)

var (
    mediaTypes = map[string]string{
        ".jpeg": "image",
        ".jpg":  "image",
        ".gif":  "image",
        ".png":  "image",
        ".mov":  "video",
        ".mp4":  "video",
        ".avi":  "video",
        ".flv":  "video",
        ".wmv":  "video",
    }
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse from body of request to get a json object.
    // fmt.Println("Received one post request")
    // decoder := json.NewDecoder(r.Body)
    // var p Post
    // if err := decoder.Decode(&p); err != nil {
    //     panic(err)

    fmt.Println("Received one upload request")

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

    if r.Method == "OPTIONS" {
        return
    }

    p := Post{      // only specify below, others empty by default
        Id: uuid.New(),
        User: r.FormValue("user"),
        Message: r.FormValue("message"),
    }

    file, header, err := r.FormFile("media_file")   // 读的content in file, other data in header
    // 现在还在uploadHandler这里，根本没送给GCS
    // media_file就是比如dog.jpeg
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }

    suffix := filepath.Ext(header.Filename) // 得到file的扩展名 比如jpg, png这些，和map比对
    // 如果是jpg/jpeg就是img, 然后还有video
    if t, ok := mediaTypes[suffix]; ok {
        p.Type = t
    } else {
        p.Type = "unknown"
    }

    err = savePost(&p, file)    // input是post and file
    if err != nil {
        http.Error(w, "Failed to save post to GCS or Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save post to GCS or Elasticsearch %v\n", err)
        return
    }

    // fmt.Fprintf(w, "Post received: %s\n", p.Message)
    fmt.Println("Post is saved successfully.")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for search")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // 这里的Header()指的是HTTP query header
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
    w.Header().Set("Content-Type", "application/json")
    if r.Method == "OPTIONS" {
        return
    }
    user := r.URL.Query().Get("user")
    keywords := r.URL.Query().Get("keywords")
    var posts []Post
    var err error
    if user != "" {
        posts, err = searchPostsByUser(user)
    } else {
        posts, err = searchPostsByKeywords(keywords)
    }
    if err != nil {
        http.Error(w, "Failed to read post from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read post from Elasticsearch %v.\n", err)
        return
    }
    js, err := json.Marshal(posts)
    if err != nil {
        http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
        fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
        return
    }
    w.Write(js)
}

