// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

// クライアントが受け取る通知の型
type Notification struct {
	// お知らせ本文
	Text string `json:"text"`
	// unixtime 形式のタイムスタンプ
	Timestamp int `json:"timestamp"`
}
