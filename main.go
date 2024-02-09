package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
	"github.com/stianeikeland/go-rpio"
        "log"
        "net"
        "strings"
    )

// HTTP POSTリクエストを送信する関数です。
func sendOpenDoorRequest(status string, minMACAddr string) {

    // 現在の時刻を取得
    currentTime := time.Now().Format(time.RFC3339)
    // POSTリクエストのボディを作成
    requestBody := []byte(fmt.Sprintf(`{"key_id": "%s", "key_status": "%s", "time": "%s"}`, minMACAddr, status, currentTime))

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

//MACアドレスの所得
func getMacAddr() (string, error) {
    ifas, err := net.Interfaces()
    if err != nil {
        return "", err
    }
    var minMACAddr string
    minMACVal := uint64(0xffffffffffff) // Initialize with maximum possible MAC value
    for _, ifa := range ifas {
        a := ifa.HardwareAddr.String()
        if a != "" {
            // Replace ":" characters to compare numerical MAC address
            a = strings.ReplaceAll(a, ":", "")
            macVal := convertMACStringToValue(a)
            if macVal < minMACVal {
                minMACVal = macVal
                minMACAddr = ifa.HardwareAddr.String()
            }
        }
    }
    return minMACAddr, nil
}

//最小のMACアドレスのみを選択 
func convertMACStringToValue(mac string) uint64 {
    var result uint64
    for i := 0; i < len(mac); i++ {
        result <<= 4
        switch {
        case '0' <= mac[i] && mac[i] <= '9':
            result |= uint64(mac[i] - '0')
        case 'a' <= mac[i] && mac[i] <= 'f':
            result |= uint64(mac[i]-'a') + 10
        case 'A' <= mac[i] && mac[i] <= 'F':
            result |= uint64(mac[i]-'A') + 10
        }
    }
    return result
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
	var keyStatus string = "CLOSE"
	var pinStatus rpio.State = TiltPin.Read()
	if  pinStatus == rpio.High {
		keyStatus = "OPEN"
	}

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

	//最小のmacアドレスの取得
	 minMACAddr, err := getMacAddr()
   	 if err != nil {
   	     log.Fatal(err)
   	 }

	// メインループ
	for {
		// GPIOピンの状態に基づいて処理（省略）
		if TiltPin.Read() == rpio.Low && pinStatus == rpio.High {
			LED("OPEN")
			fmt.Println("Tilt reset, sending OPEN status")
			sendOpenDoorRequest("OPEN", minMACAddr)
			oldTime = time.Now()
		} else if TiltPin.Read() == rpio.High && pinStatus == rpio.Low {
			LED("CLOSE")
			fmt.Println("Tilt detected, sending CLOSE status")
			sendOpenDoorRequest("CLOSE", minMACAddr)
			oldTime = time.Now()
		} else if time.Now().Sub(oldTime) > 30*time.Second && keyStatus == "OPEN" {
						   //time.Minute * 5 に変更 本番は
			fmt.Println("Tilt detected, sending OPEN ERROR status")
			sendOpenDoorRequest("Warning_Open", minMACAddr)
			oldTime = time.Now()
		}
		fmt.Println(pinStatus)
		time.Sleep(10 * time.Second) // 1秒スリープ
		//count++;
	}
}
