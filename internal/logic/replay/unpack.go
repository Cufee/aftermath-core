package replay

import (
	"archive/zip"
	"encoding/json"
	"io"
)

// World of Tanks Blitz replay.
//
// # Replay structure
//
// `*.wotbreplay` is a ZIP-archive containing:
//
// - `battle_results.dat`
// - `meta.json`
// - `data.wotreplay`

type UnpackedReplay struct {
	BattleResult battleResults `json:"results"`
	Meta         replayMeta    `json:"meta"`
}

func Unpack(file io.ReaderAt, size int64) (*UnpackedReplay, error) {
	archive, err := newZipFromReader(file, size)
	if err != nil {
		return nil, err
	}
	if len(archive.File) < 3 {
		return nil, ErrInvalidReplayFile
	}

	var data UnpackedReplay

	resultsDat, err := archive.Open("battle_results.dat")
	if err != nil {
		return nil, ErrInvalidReplayFile
	}
	result, err := decodeBattleResults(resultsDat)
	if err != nil {
		return nil, ErrInvalidReplayFile
	}
	data.BattleResult = *result

	meta, err := archive.Open("meta.json")
	if err != nil {
		return nil, ErrInvalidReplayFile
	}
	metaBytes, err := io.ReadAll(meta)
	if err != nil {
		return nil, err
	}

	return &data, json.Unmarshal(metaBytes, &data.Meta)
}

func newZipFromReader(file io.ReaderAt, size int64) (*zip.Reader, error) {
	reader, err := zip.NewReader(file, size)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
