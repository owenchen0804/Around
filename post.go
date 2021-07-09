package main

import (
    "reflect"

    "github.com/olivere/elastic/v7"
)

const (
    POST_INDEX  = "post"
)

type Post struct {
    Id      string `json:"id"`
    User    string `json:"user"`
    Message string `json:"message"`
    Url     string `json:"url"`
    Type    string `json:"type"`
}

func searchPostsByUser(user string) ([]Post, error) {
    query := elastic.NewTermQuery("user", user)
    searchResult, err := readFromES(query, POST_INDEX) // readFromES是在elasticsearch.go里
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func searchPostsByKeywords(keywords string) ([]Post, error) { // keywords表示关键字可以是多个
    query := elastic.NewMatchQuery("message", keywords)
    query.Operator("AND")   // 搜索关键字越多，结果越少
    if keywords == "" {
        query.ZeroTermsQuery("all")
    }
    searchResult, err := readFromES(query, POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []Post { // 这是个shared function
    // called by 上面两个methods
    var ptype Post
    var posts []Post    // 创建了一个没有size, cap的动态slices

    for _, item := range searchResult.Each(reflect.TypeOf(ptype)) { // 把result cast成Post type
        p := item.(Post)    // 因为返回的result里面不止是Post的信息，在item里面还含有ES的其他无关信息
        posts = append(posts, p) // 在for loop里面把p一个个的append在Posts
    }
    return posts
}

