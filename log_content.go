package logx

import (
	"bytes"
	"encoding/json"
)

// LogContent 日志输出格式
type LogContent struct {
	Prefix   string `json:"prefix"`
	Time     string `json:"time"`
	Level    string `json:"level"`
	File     string `json:"file"`
	Msg      string `json:"msg"`
	Color    bool   `json:"-"`
	LevelInt int    `json:"-"`
}

// Json LogContent to json
//	returns []byte
func (cc *LogContent) Json() []byte {
	if cc == nil {
		return []byte("")
	}
	cc.Level = levels[cc.LevelInt]
	var s bytes.Buffer
	b, _ := json.Marshal(&cc)
	if cc.Color {
		s.WriteString(logColor[cc.LevelInt])
	}
	s.Write(b)
	if cc.Color {
		s.WriteString(endColor)
	}
	return s.Bytes()
}

// Text LogContent to Text
//	returns []byte
func (cc *LogContent) Text() []byte {
	var s bytes.Buffer
	if cc.Color {
		s.WriteString(logColor[cc.LevelInt])
	}
	if cc.Prefix != "" {
		s.WriteString("[")
		s.WriteString(cc.Prefix)
		s.WriteString("]")
		s.WriteByte(' ')
	}
	s.WriteString(cc.Time)
	s.WriteByte(' ')

	s.WriteString("[")
	s.WriteString(levels[cc.LevelInt])
	s.WriteString("]")
	s.WriteByte(' ')
	s.WriteString(cc.File)
	s.WriteString(": ")
	s.WriteString(cc.Msg)

	if cc.Color {
		s.WriteString(endColor)
	}
	return s.Bytes()
}
