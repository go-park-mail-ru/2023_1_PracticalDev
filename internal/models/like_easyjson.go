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

func easyjson52421b6dDecodeGithubComGoParkMailRu20231PracticalDevInternalModels(in *jlexer.Lexer, out *Like) {
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
		case "pin_id":
			out.PinId = int(in.Int())
		case "author_id":
			out.AuthorId = int(in.Int())
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
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
func easyjson52421b6dEncodeGithubComGoParkMailRu20231PracticalDevInternalModels(out *jwriter.Writer, in Like) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"pin_id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.PinId))
	}
	{
		const prefix string = ",\"author_id\":"
		out.RawString(prefix)
		out.Int(int(in.AuthorId))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Like) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson52421b6dEncodeGithubComGoParkMailRu20231PracticalDevInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Like) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson52421b6dEncodeGithubComGoParkMailRu20231PracticalDevInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Like) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson52421b6dDecodeGithubComGoParkMailRu20231PracticalDevInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Like) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson52421b6dDecodeGithubComGoParkMailRu20231PracticalDevInternalModels(l, v)
}
