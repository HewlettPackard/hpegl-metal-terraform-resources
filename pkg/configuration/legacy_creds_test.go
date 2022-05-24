// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package configuration

import (
	"bytes"
	"testing"
)

const (
	fakeQJWT = `rest_url: http://15.242.208.109
original_url: http://15.242.208.109
user: h1@quattronetworks.com
jwt: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJFTk
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
