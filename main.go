package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/stianeikeland/go-rpio"
)

// HTTP POSTリクエストを送信する関数です。
func sendOpenDoorRequest(status string) {
    // 現在の時刻を取得
    currentTime := time.Now().Format(time.RFC3339)
    
    // POSTリクエストのボディを作成
    requestBody := []byte(fmt.Sprintf(`{"key_status": "%s", "time": "%s"}`, status, currentTime))

    // http.Post を使用してリクエストを送信
    response, err := http.Post("https://golang-line-api-tokazaki-20240207.onrender.com/open_door", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        fmt.Printf("POSTリクエストの送信に失敗しました: %v\n", err)
        return
    }
    defer response.Body.Close()

    // ステータスコードを確認
    fmt.Printf("Response Status: %s\n", response.Status)
}

func main() {
	// Raspberry Pi GPIO初期化
	if err := rpio.Open(); err != nil {
		fmt.Println("Error initializing GPIO:", err)
		return
	}
	defer rpio.Close()

	// GPIOピンの設定
	TiltPin := rpio.Pin(17) // GPIOピン17を使用
	TiltPin.Input()

	// LEDの現在の状態を保持する変数
	var keyStatus string = "OPEN"
	var pinStatus rpio.State = rpio.High

	// LED関数（LEDの制御）
	LED := func(status string) {
		keyStatus = status
		if keyStatus == "OPEN" {
			pinStatus = rpio.Low
		} else if keyStatus == "CLOSE" {
			pinStatus = rpio.High
		}

		Gpin := rpio.Pin(2) // GPIOピン2を使用
		Rpin := rpio.Pin(3) // GPIOピン3を使用

		// ピンを出力モードに設定
		Gpin.Output()
		Rpin.Output()

		if keyStatus == "CLOSE" {
			Rpin.Low()
			Gpin.High()
		} else if keyStatus == "OPEN" {
			Rpin.High()
			Gpin.Low()
		} else {
			fmt.Println("LED ERROR")
		}
	}

	var oldTime time.Time = time.Now()
	//var int coutnt

	// メインループ
	for {
		// GPIOピンの状態に基づいて処理（省略）
		if TiltPin.Read() == rpio.Low && pinStatus == rpio.High {
			LED("OPEN")
			fmt.Println("Tilt reset, sending OPEN status")
			//sendOpenDoorRequest("OPEN")
			oldTime = time.Now()
		} else if TiltPin.Read() == rpio.High && pinStatus == rpio.Low {
			LED("CLOSE")
			fmt.Println("Tilt detected, sending CLOSE status")
			//sendOpenDoorRequest("CLOSE")
			oldTime = time.Now()
		} else if time.Now().Sub(oldTime) > 5*time.Second && keyStatus == "OPEN" {
						   //time.Minute * 5 に変更 本番は
			fmt.Println("Tilt detected, sending OPEN ERROR status")
			//sendOpenDoorRequest("Warning_Open")
			oldTime = time.Now()
		}
		fmt.Println(pinStatus)
		time.Sleep(1 * time.Second) // 1秒スリープ
		//count++;
	}
}
