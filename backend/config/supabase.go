package config

import (
	"log"
	"os"

	"github.com/supabase-community/supabase-go"
)

var SupabaseClient *supabase.Client

// InitSupabase Supabaseクライアントを初期化
func InitSupabase() {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URLとSUPABASE_KEYの設定が必要です。")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, &supabase.ClientOptions{})
	if err != nil {
		log.Fatalf("Supabaseの初期化に失敗:%v", err)
	}
	SupabaseClient = client
	log.Println("Supabaseクライアント初期化完了。")
}
