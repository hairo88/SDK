package main

import (
	"bytes"
	"fmt"
	"time"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stianeikeland/go-rpio"
)

// HTTP POSTリクエストを送信する関数です。
func sendOpenDoorRequest(keyStatus string) {
    // POSTリクエストのボディを作成
    requestBody := []byte(fmt.Sprintf(`{"key_status": "%s"}`, keyStatus))

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
	// AWSセッションを作成
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)

	// S3サービスクライアントを作成
	s3Client := s3.New(sess)

	// バケット名とアップロードするファイル名
	bucketName := "magickeybucket" // 実際のバケット名に変更してください
	fileName := "current-time.txt"

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
	var ledState string = "GREEN"

	var keyState rpio.State = rpio.High

	// LED関数（LEDの制御）
	LED := func(color string) {
		ledState = color

		if color == "GREEN" {
			keyState = rpio.High
		} else if color == "RED" {
			keyState = rpio.Low
		}

		Gpin := rpio.Pin(2) // GPIOピン2を使用
		Rpin := rpio.Pin(3) // GPIOピン3を使用

		// ピンを出力モードに設定
		Gpin.Output()
		Rpin.Output()

		if color == "RED" {
			Rpin.High()
			Gpin.Low()
		} else if color == "GREEN" {
			Rpin.Low()
			Gpin.High()
		} else {
			fmt.Println("LED ERROR")
		}
	}

	// 初期状態としてLEDを緑色に設定
	LED("GREEN")

	// メインループ
	for {
			if TiltPin.Read() == rpio.Low && keyState == rpio.High {
				LED("RED")
				fmt.Println("Tilt!")

				// 現在の時間を取得
				currentTime := time.Now().Format(time.RFC3339)
				content := fmt.Sprintf("LED State: %s, Time: %s", ledState, currentTime)

				// ファイルの内容をバイト配列に変換
				fileContent := []byte(content)

				// PutObject入力を作成
				putObjectInput := &s3.PutObjectInput{
					Bucket:        aws.String(bucketName),
					Key:           aws.String(fileName),
					Body:          bytes.NewReader(fileContent),
					ContentLength: aws.Int64(int64(len(fileContent))),
					ContentType:   aws.String("text/plain"),
				}

				// S3バケットにファイルをアップロード
				_, err := s3Client.PutObject(putObjectInput)
				if err != nil {
					fmt.Printf("Unable to upload %q to %q, %v", fileName, bucketName, err)
				} else {
					fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucketName)
				}
		} else if TiltPin.Read() == rpio.High && keyState == rpio.Low {
				LED("GREEN")
				fmt.Println("Tilt!")

				// 現在の時間を取得
				currentTime := time.Now().Format(time.RFC3339)
				content := fmt.Sprintf("LED State: %s, Time: %s", ledState, currentTime)

				// ファイルの内容をバイト配列に変換
				fileContent := []byte(content)

				sendOpenDoorRequest("OPEN")

				// PutObject入力を作成
				putObjectInput := &s3.PutObjectInput{
					Bucket:        aws.String(bucketName),
					Key:           aws.String(fileName),
					Body:          bytes.NewReader(fileContent),
					ContentLength: aws.Int64(int64(len(fileContent))),
					ContentType:   aws.String("text/plain"),
				}

				// S3バケットにファイルをアップロード
				_, err := s3Client.PutObject(putObjectInput)
				if err != nil {
					fmt.Printf("Unable to upload %q to %q, %v", fileName, bucketName, err)
				} else {
					fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucketName)
				}
		}
	}
}
