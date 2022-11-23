package logx

import (
	"bytes"
	"encoding/json"
)

// LogContent 日志输出格式
type LogContent struct {
	Prefix  string `json:"prefix"`
	Level   int    `json:"level"`
	Package string `json:"package"`
	File    string `json:"file"`
	Msg     string `json:"msg"`
	Time    string `json:"time"`
	Color   bool   `json:"-"`
}

// Json LogContent to json
//	returns []byte
func (cc *LogContent) Json() []byte {
	if cc == nil {
		return []byte("")
	}
	var s bytes.Buffer
	b, _ := json.Marshal(&cc)
	if cc.Color {
		s.WriteString(logColor[cc.Level])
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
		s.WriteString(logColor[cc.Level])
	}
	if cc.Prefix != "" {
		s.WriteString("[")
		s.WriteString(cc.Prefix)
		s.WriteString("]")
		s.WriteByte(' ')
	}
	s.WriteString(cc.Time)
	s.WriteByte(' ')
	s.WriteString(levels[cc.Level])
	s.WriteByte(' ')
	s.WriteString(cc.File)
	s.WriteString(": ")
	s.WriteString(cc.Msg)

	if cc.Color {
		s.WriteString(endColor)
	}
	return s.Bytes()
}
