// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonA818f49aDecodeGithubComGoParkMailRu20231PracticalDevInternalModels(in *jlexer.Lexer, out *Session) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.UserId = int(in.Int())
		case "email":
			out.UserEmail = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonA818f49aEncodeGithubComGoParkMailRu20231PracticalDevInternalModels(out *jwriter.Writer, in Session) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.UserId))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.UserEmail))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Session) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA818f49aEncodeGithubComGoParkMailRu20231PracticalDevInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Session) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA818f49aEncodeGithubComGoParkMailRu20231PracticalDevInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Session) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA818f49aDecodeGithubComGoParkMailRu20231PracticalDevInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Session) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA818f49aDecodeGithubComGoParkMailRu20231PracticalDevInternalModels(l, v)
}
