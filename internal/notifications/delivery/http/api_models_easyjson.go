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

func easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalNotificationsDeliveryHttp(in *jlexer.Lexer, out *listResponse) {
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
				out.Notifications = nil
			} else {
				in.Delim('[')
				if out.Notifications == nil {
					if !in.IsDelim(']') {
						out.Notifications = make([]models.Notification, 0, 0)
					} else {
						out.Notifications = []models.Notification{}
					}
				} else {
					out.Notifications = (out.Notifications)[:0]
				}
				for !in.IsDelim(']') {
					var v1 models.Notification
					(v1).UnmarshalEasyJSON(in)
					out.Notifications = append(out.Notifications, v1)
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
func easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalNotificationsDeliveryHttp(out *jwriter.Writer, in listResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"items\":"
		out.RawString(prefix[1:])
		if in.Notifications == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Notifications {
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
	easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalNotificationsDeliveryHttp(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v listResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC0ea9389EncodeGithubComGoParkMailRu20231PracticalDevInternalNotificationsDeliveryHttp(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *listResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalNotificationsDeliveryHttp(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *listResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC0ea9389DecodeGithubComGoParkMailRu20231PracticalDevInternalNotificationsDeliveryHttp(l, v)
}
