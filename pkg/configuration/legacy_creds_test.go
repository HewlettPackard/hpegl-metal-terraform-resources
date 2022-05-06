// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package configuration

import (
	"bytes"
	"testing"
)

const (
	fakeQJWT = `rest_url: http://15.242.208.109:902
original_url: http://15.242.208.109:902
user: h1@quattronetworks.com
jwt: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJFTkRORFU1T1VFM05EZzNRVEkxUVVRd01qRTBOVFE0UkVaQk9VVTBRVGd5TWtJek5URXhRZyJ9.eyJuaWNrbmFtZSI6ImgxIiwibmFtZSI6ImgxQHF1YXR0cm9uZXR3b3Jrcy5jb20iLCJwaWN0dXJlIjoiaHR0cHM6Ly9zLmdyYXZhdGFyLmNvbS9hdmF0YXIvMDE
zZWQ3ZmMzOTZkZGNhN2RkYzI4YWQyMzgwN2NlZjk_cz00ODAmcj1wZyZkPWh0dHBzJTNBJTJGJTJGY2RuLmF1dGgwLmNvbSUyRmF2YXRhcnMlMkZoMS5wbmciLCJ1cGRhdGVkX2F0IjoiMjAyMC0wNC0wMlQwMToyMzoyMC42NjZaIiwiZW1haWwiOiJoMUBxdWF0dHJvbmV0d29ya3MuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUs
ImlzcyI6Imh0dHBzOi8vcXVhdHRyb3Rlc3QuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfDViOTY5YjRiZjE2NDI3MjFhZDg0YzI1MyIsImF1ZCI6IklBYlVDQnEwVHp5YTBxQ3o3UVU3blFkQnFYbTlTRHY2IiwiaWF0IjoxNTg1NzkwNjAwLCJleHAiOjE1ODU4MjY2MDB9.FlVS-LdPw03XSHqULM8QKNulILFCXbXLtjdOZo4vygHJ
opJUXGhOJu0oGkiXeG5x1KCCv1vg2dGsCg6I_sgEy_Nek5ASG8VdvR6JalnNXWFMpE_wnK7RTIIFELMKH4WmxJBaYEFmThQ9y-yggWitmWb6Mgs3uBEoe_d5wsqEpdVh8Zc6h0qft4-Vdf3P_wdBvSvUpaFVlqJzIPiJWGMkdBik0kwzZk1xi-m8vs2T3v-zqGKSgbAM5V7xxmkfGdywc9MUM-fBWxbfzS2nyGxU9vB9tQaBNfZyW8-WM
oPoB38oZEA-y7JC2sEBOS7P7NnjuJefLffRLtFjF3nvK--84Q
member_id: B901063C-DB35-4FF8-8A78-EA44C62C61F8
no_tls: true
`
)

func TestGetConfig(t *testing.T) {
	q, err := parseStream(bytes.NewBufferString(fakeQJWT))
	if err != nil {
		t.Fatal((err))
	}
	if q.RestURL == "" {
		t.Fatal("RestURL empty")
	}
	if q.OriginalURL == "" {
		t.Fatal("OriginalURL empty")
	}
	if q.Token == "" {
		t.Fatal("Token empty")
	}
	if q.MemberID == "" {
		t.Fatal("MemberID empty")
	}
}
