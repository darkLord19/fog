package ai

import "testing"

func TestStreamJSONParserExtractsTextAndConversationID(t *testing.T) {
	var chunks []string
	parser := newStreamJSONParser(func(chunk string) {
		chunks = append(chunks, chunk)
	})

	parser.Feed([]byte(`{"type":"assistant","text":"hello"}` + "\n"))
	parser.Feed([]byte(`{"type":"result","session_id":"sess-123"}` + "\n"))
	parser.Close()

	if parser.Output() != "hello" {
		t.Fatalf("unexpected parser output: %q", parser.Output())
	}
	if parser.ConversationID() != "sess-123" {
		t.Fatalf("unexpected conversation id: %q", parser.ConversationID())
	}
	if len(chunks) != 1 || chunks[0] != "hello" {
		t.Fatalf("unexpected streamed chunks: %+v", chunks)
	}
}

func TestStreamJSONParserIgnoresUserRole(t *testing.T) {
	parser := newStreamJSONParser(nil)
	parser.Feed([]byte(`{"role":"user","text":"ignore this"}` + "\n"))
	parser.Feed([]byte(`{"role":"assistant","text":"keep this"}` + "\n"))
	parser.Close()

	if parser.Output() != "keep this" {
		t.Fatalf("unexpected parser output: %q", parser.Output())
	}
}

func TestStreamJSONParserPreservesLeadingSpacesAcrossChunks(t *testing.T) {
	var chunks []string
	parser := newStreamJSONParser(func(chunk string) {
		chunks = append(chunks, chunk)
	})

	parser.Feed([]byte(`{"type":"assistant","text":"Hello"}` + "\n"))
	parser.Feed([]byte(`{"type":"assistant","text":" world"}` + "\n"))
	parser.Close()

	if parser.Output() != "Hello world" {
		t.Fatalf("unexpected parser output: %q", parser.Output())
	}
	if len(chunks) != 2 || chunks[0] != "Hello" || chunks[1] != " world" {
		t.Fatalf("unexpected streamed chunks: %+v", chunks)
	}
}
