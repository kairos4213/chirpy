package main

import "testing"

func TestCheckForProfaneWords(t *testing.T) {
	bodyString := "I hear Mastodon is better than Chirpy. sharbert I need to migrate"
	want := "I hear Mastodon is better than Chirpy. **** I need to migrate"
	result := checkForProfaneWords(bodyString)
	if want != result {
		t.Fatalf(`checkForProfaneWords(%s) = %s, want %s`, bodyString, result, want)
	}
}
