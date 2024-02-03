package external

import (
	"testing"
)

func TestGetReplayData(t *testing.T) {
	url := "https://cdn.discordapp.com/attachments/1088461326432608256/1203465847390347294/1706997697_karelia_8fc8d9abc7af7113df394b0b0a02dcc8.wotbreplay"
	replay, err := GetReplayData(url)
	if err != nil {
		t.Fatal(err)
	}

	if replay == nil {
		t.Fatal("replay is nil")
	}
}
