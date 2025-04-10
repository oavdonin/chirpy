package main

import (
	"fmt"
	"testing"
)

func TestReplaceProfaneWords(t *testing.T) {
	var tests = []struct {
		given RequestBody
		want  RequestBody
	}{
		{
			given: RequestBody{Body: "I had something interesting for breakfast"},
			want:  RequestBody{Body: "I had something interesting for breakfast"},
		},
		{
			given: RequestBody{Body: "I hear Mastodon is better than Chirpy. sharbert I need to migrate"},
			want:  RequestBody{Body: "I hear Mastodon is better than Chirpy. **** I need to migrate"},
		},
		{
			given: RequestBody{Body: "I really need a kerfuffle to go to bed sooner, Fornax !"},
			want:  RequestBody{Body: "I really need a **** to go to bed sooner, **** !"},
		},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("Test #%d", i)
		t.Run(testname, func(t *testing.T) {
			replaceProfaneWords(&tt.given)
			if tt.given != tt.want {
				t.Errorf("got %s, want %s", tt.given.Body, tt.want.Body)
			}
		})
	}
}
