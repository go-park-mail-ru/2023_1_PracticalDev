// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package http

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

func easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalShortenerDeliveryHttp(in *jlexer.Lexer, out *url) {
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
		case "url":
			out.URL = string(in.String())
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
func easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalShortenerDeliveryHttp(out *jwriter.Writer, in url) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix[1:])
		out.String(string(in.URL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v url) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalShortenerDeliveryHttp(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v url) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalShortenerDeliveryHttp(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *url) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalShortenerDeliveryHttp(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *url) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalShortenerDeliveryHttp(l, v)
}