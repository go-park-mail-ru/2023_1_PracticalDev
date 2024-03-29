// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package http

import (
	json "encoding/json"
	models "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
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

func easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp(in *jlexer.Lexer, out *listResponse) {
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
		case "items":
			if in.IsNull() {
				in.Skip()
				out.Items = nil
			} else {
				in.Delim('[')
				if out.Items == nil {
					if !in.IsDelim(']') {
						out.Items = make([]models.Comment, 0, 1)
					} else {
						out.Items = []models.Comment{}
					}
				} else {
					out.Items = (out.Items)[:0]
				}
				for !in.IsDelim(']') {
					var v1 models.Comment
					(v1).UnmarshalEasyJSON(in)
					out.Items = append(out.Items, v1)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp(out *jwriter.Writer, in listResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"items\":"
		out.RawString(prefix[1:])
		if in.Items == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Items {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v listResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v listResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *listResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *listResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp(l, v)
}
func easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp1(in *jlexer.Lexer, out *createResponse) {
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
			out.ID = int(in.Int())
		case "author_id":
			out.AuthorID = int(in.Int())
		case "pin_id":
			out.PinID = int(in.Int())
		case "text":
			out.Text = string(in.String())
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
func easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp1(out *jwriter.Writer, in createResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"author_id\":"
		out.RawString(prefix)
		out.Int(int(in.AuthorID))
	}
	{
		const prefix string = ",\"pin_id\":"
		out.RawString(prefix)
		out.Int(int(in.PinID))
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v createResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v createResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *createResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *createResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp1(l, v)
}
func easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp2(in *jlexer.Lexer, out *createRequest) {
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
		case "text":
			out.Text = string(in.String())
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
func easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp2(out *jwriter.Writer, in createRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix[1:])
		out.String(string(in.Text))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v createRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v createRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *createRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *createRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalCommentsDeliveryHttp2(l, v)
}
