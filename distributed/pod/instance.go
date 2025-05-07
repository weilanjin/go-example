package pod

import (
	"encoding/json"
	"time"
)

type Instance struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Version  string    `json:"version"`
	Address  string    `json:"address"`
	Port     int       `json:"port"`
	LashBeat time.Time `json:"lash_beat"` // 最后一次心跳时间
}

func (p *Instance) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Instance) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
