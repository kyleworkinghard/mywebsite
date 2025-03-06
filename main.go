// main.go
package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

const GROK_API_KEY = "xai-PwA5Wh7pOYnJPcovkOceOi7y35zNSNW5pRFO2COxCj0AR5cTH3OlfthU1oe5ZqYfsp1LrMb01i14I3ad"

type GrokRequest struct {
	Message string `json:"message"`
}

type GrokResponse struct {
	Response string `json:"response"`
}

func handleGrokChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GrokRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	grokReq := map[string]interface{}{
		"model": "grok-beta",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": req.Message,
			},
		},
	}

	jsonData, err := json.Marshal(grokReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	request, err := http.NewRequest("POST", "https://api.x.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+GROK_API_KEY)

	resp, err := client.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var grokResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&grokResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 从 Grok 响应中提取文本
	choices := grokResp["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := message["content"].(string)

	response := GrokResponse{
		Response: content,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/chat.html"))
	tmpl.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// 如果不是首页路径，返回404
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	data := map[string]interface{}{
		"Title":   "我的网站",
		"Message": "下面是基于Grok2的聊天机器人",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("模板渲染错误: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	// 设置静态文件服务
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 设置路由
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/api/grok-chat", handleGrokChat)

	// 启动服务器
	log.Println("服务器启动在 http://localhost:3389")
	if err := http.ListenAndServe(":3389", nil); err != nil {
		log.Fatal(err)
	}
}
