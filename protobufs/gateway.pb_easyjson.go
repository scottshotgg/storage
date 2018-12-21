// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package protobuf

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

func easyjsonC47c71faDecodeGithubComScottshotggStorageProtobufs(in *jlexer.Lexer, out *storageClient) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
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
func easyjsonC47c71faEncodeGithubComScottshotggStorageProtobufs(out *jwriter.Writer, in storageClient) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v storageClient) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC47c71faEncodeGithubComScottshotggStorageProtobufs(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v storageClient) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC47c71faEncodeGithubComScottshotggStorageProtobufs(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *storageClient) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC47c71faDecodeGithubComScottshotggStorageProtobufs(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *storageClient) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC47c71faDecodeGithubComScottshotggStorageProtobufs(l, v)
}
func easyjsonC47c71faDecodeGithubComScottshotggStorageProtobufs1(in *jlexer.Lexer, out *Changelog) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ID":
			out.ID = string(in.String())
		case "timestamp":
			out.Timestamp = int64(in.Int64())
		case "itemID":
			out.ItemID = string(in.String())
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
func easyjsonC47c71faEncodeGithubComScottshotggStorageProtobufs1(out *jwriter.Writer, in Changelog) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"ID\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ID))
	}
	if in.Timestamp != 0 {
		const prefix string = ",\"timestamp\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Timestamp))
	}
	if in.ItemID != "" {
		const prefix string = ",\"itemID\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ItemID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Changelog) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC47c71faEncodeGithubComScottshotggStorageProtobufs1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Changelog) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC47c71faEncodeGithubComScottshotggStorageProtobufs1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Changelog) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC47c71faDecodeGithubComScottshotggStorageProtobufs1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Changelog) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC47c71faDecodeGithubComScottshotggStorageProtobufs1(l, v)
}
