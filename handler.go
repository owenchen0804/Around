package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse from body of request to get a json object.
    fmt.Println("Received one post request")

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") // CROS handling??

    if r.Method == "OPTIONS" { // ??
        return
    }

    decoder := json.NewDecoder(r.Body)
    var p Post
    if err := decoder.Decode(&p); err != nil { // &p ?? 
        panic(err)
    }

    fmt.Fprintf(w, "Post received: %s\n", p.Message)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for search")

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
    w.Header().Set("Content-Type", "application/json")

    if r.Method == "OPTIONS" {
        return
    }

    user := r.URL.Query().Get("user") // 像Servlet里面去拿query的parameter
    keywords := r.URL.Query().Get("keywords") // 类似getParam 

    var posts []Post
    var err error
    if user != "" { // empty string 相当于 zero value
        posts, err = searchPostsByUser(user)    // 有user的话就按照user搜索，没有user才按照keywords
    } else {
        posts, err = searchPostsByKeywords(keywords)
    }

    if err != nil {
        http.Error(w, "Failed to read post from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read post from Elasticsearch %v.\n", err)
        return
    }

    js, err := json.Marshal(posts)  // js就是json格式的result, Marshal == marsh 意思是deserialize结果
                                    // 得到json string
    if err != nil {
        http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
        fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
        return
    }
    w.Write(js)
}
