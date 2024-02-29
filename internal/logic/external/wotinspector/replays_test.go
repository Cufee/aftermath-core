package wotinspector

import (
	"encoding/json"
	"testing"
)

func TestGetReplayData(t *testing.T) {
	url := "https://cdn.discordapp.com/attachments/1088461326432608256/1204786498012971058/1706997697_karelia_8fc8d9abc7af7113df394b0b0a02dcc81.wotbreplay?ex=65d5ffdc&is=65c38adc&hm=44a3880be552cf8d0e7aa396a832240af5986d94fc079c79df571315c3c906fe&"
	replay, err := GetReplayData(url)
	if err != nil {
		t.Fatal(err)
	}

	if replay == nil {
		t.Fatal("replay is nil")
	}

	data, _ := json.Marshal(replay)
	t.Log(string(data))
}
