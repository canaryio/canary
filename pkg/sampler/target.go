package sampler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

type Target struct {
	URL      string
	Name     string
	Interval int
	// metadata
	Tags               []string
	Attributes         map[string]string
	Hash               string
	RequestHeaders     map[string]string
	InsecureSkipVerify bool
}

func (t *Target) SetHash() {
	jsonTarget, _ := json.Marshal(t)
	hasher := md5.New()
	hasher.Write(jsonTarget)
	t.Hash = hex.EncodeToString(hasher.Sum(nil))
}
