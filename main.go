// main.go
package main

import (
    "html/template"
    "log"
    "net/http"
)

func main() {
    // 设置静态文件服务
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // 设置路由
    http.HandleFunc("/", homeHandler)

    // 启动服务器
    log.Println("服务器启动在 http://localhost:3389")
    if err := http.ListenAndServe(":3389", nil); err != nil {
        log.Fatal(err)
    }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    // 如果不是首页路径，返回404
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    tmpl := template.Must(template.ParseFiles("templates/home.html"))
    data := map[string]interface{}{
        "Title": "我的网站",
        "Message": "欢迎访问！",
    }
    
    if err := tmpl.Execute(w, data); err != nil {
        log.Printf("模板渲染错误: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
