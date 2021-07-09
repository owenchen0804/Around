package main

import (
    "context"

    "github.com/olivere/elastic/v7"
)

const (
        ES_URL = "http://10.128.0.2:9200" // external IP可能会变，且这个Internal允许所有IP地址访问
)

func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    client, err := elastic.NewClient(
        elastic.SetURL(ES_URL),
        elastic.SetBasicAuth("elastic", "12345678")) // 实际生产要用环境变量来存密码
    if err != nil {
        return nil, err
    }

    searchResult, err := client.Search(). // code能到这里说明没有error
        Index(index).   // index具体是啥也是Input
        Query(query).   // query是从外面传进来的
        Pretty(true).
        Do(context.Background()) // 这里的"."做法叫fluent API 相当于连环call
    if err != nil {
        return nil, err
    }

    return searchResult, nil
}