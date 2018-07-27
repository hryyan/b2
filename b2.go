// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package b2 is a go library for backblaze B2 Cloud Storage.
package b2

// B2 is used to initialize your b2 account and applicationkey
type B2 struct {
	AccountId      string
	ApplicationKey string
	auth           AuthResponse
}

func (b *B2) GetAuth() AuthResponse {
	return b.auth
}

func (b *B2) SetAuth(auth AuthResponse) {
	b.auth = auth
}
