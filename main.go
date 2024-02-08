// package main

// import (
// 	"bytes"
// 	"fmt"
// 	"time"
// 	"net/http"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/s3"
// 	"github.com/stianeikeland/go-rpio"
// )

// // HTTP POSTリクエストを送信する関数です。
// func sendOpenDoorRequest(keyStatus string) {
//     // POSTリクエストのボディを作成
//     requestBody := []byte(fmt.Sprintf(`{"key_status": "%s"}`, keyStatus))

//     // http.Post を使用してリクエストを送信
//     response, err := http.Post("https://golang-line-api-tokazaki-20240207.onrender.com/open_door", "application/json", bytes.NewBuffer(requestBody))
//     if err != nil {
//         fmt.Printf("POSTリクエストの送信に失敗しました: %v\n", err)
//         return
//     }
//     defer response.Body.Close()

//     // ステータスコードを確認
//     fmt.Printf("Response Status: %s\n", response.Status)


// }

// func main() {
// 	// // AWSセッションを作成
// 	// sess, _ := session.NewSession(&aws.Config{
// 	// 	Region: aws.String("ap-northeast-1")},
// 	// )

// 	// // S3サービスクライアントを作成
// 	// s3Client := s3.New(sess)

// 	// // バケット名とアップロードするファイル名
// 	// bucketName := "magickeybucket" // 実際のバケット名に変更してください
// 	// fileName := "current-time.txt"

// 	// Raspberry Pi GPIO初期化
// 	if err := rpio.Open(); err != nil {
// 		fmt.Println("Error initializing GPIO:", err)
// 		return
// 	}
// 	defer rpio.Close()

// 	// GPIOピンの設定
// 	TiltPin := rpio.Pin(17) // GPIOピン17を使用
// 	TiltPin.Input()

// 	// LEDの現在の状態を保持する変数
// 	var ledState string = "OPEN"

// 	var keyState rpio.State = rpio.High

// 	// LED関数（LEDの制御）
// 	LED := func(color string) {
// 		ledState = color

// 		if color == "OPEN" {
// 			keyState = rpio.High
// 		} else if color == "CLOSE" {
// 			keyState = rpio.Low
// 		}

// 		Gpin := rpio.Pin(2) // GPIOピン2を使用
// 		Rpin := rpio.Pin(3) // GPIOピン3を使用

// 		// ピンを出力モードに設定
// 		Gpin.Output()
// 		Rpin.Output()

// 		if color == "CLOSE" {
// 			Rpin.High()
// 			Gpin.Low()
// 		} else if color == "OPEN" {
// 			Rpin.Low()
// 			Gpin.High()
// 		} else {
// 			fmt.Println("LED ERROR")
// 		}
// 	}

// 	// 初期状態としてLEDを緑色に設定
// 	LED("OPEN")

// 	// メインループ
// 	for {
// 			if TiltPin.Read() == rpio.Low && keyState == rpio.High {
// 				LED("CLOSE")
// 				fmt.Println("Tilt!")

// 				// 現在の時間を取得
// 				currentTime := time.Now().Format(time.RFC3339)
// 				content := fmt.Sprintf("LED State: %s, Time: %s", ledState, currentTime)

// 				// // ファイルの内容をバイト配列に変換
// 				// fileContent := []byte(content)

// 				// // PutObject入力を作成
// 				// putObjectInput := &s3.PutObjectInput{
// 				// 	Bucket:        aws.String(bucketName),
// 				// 	Key:           aws.String(fileName),
// 				// 	Body:          bytes.NewReader(fileContent),
// 				// 	ContentLength: aws.Int64(int64(len(fileContent))),
// 				// 	ContentType:   aws.String("text/plain"),
// 				// }

// 				// // S3バケットにファイルをアップロード
// 				// _, err := s3Client.PutObject(putObjectInput)
// 				// if err != nil {
// 				// 	fmt.Printf("Unable to upload %q to %q, %v", fileName, bucketName, err)
// 				// } else {
// 				// 	fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucketName)
// 				// }
// 		} else if TiltPin.Read() == rpio.High && keyState == rpio.Low {
// 				LED("OPEN")
// 				fmt.Println("Tilt!")

// 				// 現在の時間を取得
// 				currentTime := time.Now().Format(time.RFC3339)
// 				content := fmt.Sprintf("LED State: %s, Time: %s", ledState, currentTime)

// 				// // ファイルの内容をバイト配列に変換
// 				// fileContent := []byte(content)

// 		// 		sendOpenDoorRequest("OPEN")

// 		// 		// PutObject入力を作成
// 		// 		putObjectInput := &s3.PutObjectInput{
// 		// 			Bucket:        aws.String(bucketName),
// 		// 			Key:           aws.String(fileName),
// 		// 			Body:          bytes.NewReader(fileContent),
// 		// 			ContentLength: aws.Int64(int64(len(fileContent))),
// 		// 			ContentType:   aws.String("text/plain"),
// 		// 		}

// 		// 		// S3バケットにファイルをアップロード
// 		// 		_, err := s3Client.PutObject(putObjectInput)
// 		// 		if err != nil {
// 		// 			fmt.Printf("Unable to upload %q to %q, %v", fileName, bucketName, err)
// 		// 		} else {
// 		// 			fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucketName)
// 		// 		}
// 		}

// 	}
// }



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
			pinStatus = rpio.High
		} else if keyStatus == "CLOSE" {
			pinStatus = rpio.Low
		}

		Gpin := rpio.Pin(2) // GPIOピン2を使用
		Rpin := rpio.Pin(3) // GPIOピン3を使用

		// ピンを出力モードに設定
		Gpin.Output()
		Rpin.Output()

		if keyStatus == "CLOSE" {
			Rpin.High()
			Gpin.Low()
		} else if keyStatus == "OPEN" {
			Rpin.Low()
			Gpin.High()
		} else {
			fmt.Println("LED ERROR")
		}
	}

	// メインループ
	for {
		// GPIOピンの状態に基づいて処理（省略）
		if TiltPin.Read() == rpio.Low && pinStatus == rpio.High {
			LED("CLOSE")
			fmt.Println("Tilt detected, sending CLOSE status")
			sendOpenDoorRequest("CLOSE")
		} else if TiltPin.Read() == rpio.High && pinStatus == rpio.Low {
			LED("OPEN")
			fmt.Println("Tilt reset, sending OPEN status")
			sendOpenDoorRequest("OPEN")
		}
	}
}
