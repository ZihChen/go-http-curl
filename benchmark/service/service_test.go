package service

import (
	"testing"
	"time"
)

// 驗證ConnectPool在多線程請求結果是否正確
// 測試重點：
// 1.各個Server端取得的query params及headers需與線程請求帶的參數一致
// 2.各個Server端顯示的remote address的port號是否都來自同一個
// 3.每個線程只會對Server建立一次連線，Server關閉後顯示的reconnect count需等於0
func TestHttpConnectPool(t *testing.T) {
	s := NewService()

	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(1 * time.Second)
			}
			go func() {
				s.MethodA()
			}()
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(1 * time.Second)
			}

			go func() {
				s.MethodB()
			}()
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(1 * time.Second)
			}

			go func() {
				s.MethodC()
			}()
		}
	}()

	time.Sleep(15 * time.Second)
}

// 驗證普通連線(非ConnectPool)
func TestHttpCommonConnect(t *testing.T) {
	s := NewService()

	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(1 * time.Second)
			}
			go func() {
				s.MethodD()
			}()
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(1 * time.Second)
			}

			go func() {
				s.MethodF()
			}()
		}
	}()

	time.Sleep(10 * time.Second)
}
