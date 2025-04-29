package main

import (
	"log"
	"net/http"
	"os"

	"cloudtech_forum/handler"
	"cloudtech_forum/repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// .envファイルから環境変数を読み込む
	godotenv.Load()
}

// CORSミドルウェア
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORSヘッダーを設定（必要に応じて制限することも可能）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// OPTIONSメソッドへの即時レスポンス（Preflightリクエスト対応）
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 通常処理
		next.ServeHTTP(w, r)
	})
}

func main() {

	// 環境変数からデータを取得
	apiport := os.Getenv("API_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// データベース接続
	err := repository.InitDB(username, password, host, port, dbname)
	if err != nil {
		log.Fatalf("データベースに接続できません: %v", err)
	}
	defer repository.CloseDB()

	// ルーター定義
	r := mux.NewRouter()
	r.HandleFunc("/posts", handler.Create).Methods("POST")
	r.HandleFunc("/posts", handler.Index).Methods("GET")
	r.HandleFunc("/posts/{id:[0-9]+}", handler.Show).Methods("GET")
	r.HandleFunc("/posts/{id:[0-9]+}", handler.Update).Methods("PUT")
	r.HandleFunc("/posts/{id:[0-9]+}", handler.Delete).Methods("DELETE")

	// CORSミドルウェアを適用
	corsRouter := enableCORS(r)

	// APIサーバ起動
	log.Println("APIサーバを起動しました。ポート: " + apiport)
	if err := http.ListenAndServe(":"+apiport, corsRouter); err != nil {
		log.Fatal(err)
	}
}
