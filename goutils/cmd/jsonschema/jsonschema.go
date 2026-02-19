// jsonschema - infer JSON schema from a JSON document or JSONL sample
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type schema struct {
	Type       string             `json:"type"`
	Properties map[string]*schema `json:"properties,omitempty"`
	Items      *schema            `json:"items,omitempty"`
	Required   []string           `json:"required,omitempty"`
	Examples   []interface{}      `json:"examples,omitempty"`
}

func infer(v interface{}) *schema {
	switch t := v.(type) {
	case map[string]interface{}:
		s := &schema{Type: "object", Properties: map[string]*schema{}}
		keys := make([]string, 0, len(t)); for k := range t { keys = append(keys, k) }
		sort.Strings(keys)
		for _, k := range keys { s.Properties[k] = infer(t[k]); s.Required = append(s.Required, k) }
		return s
	case []interface{}:
		s := &schema{Type: "array"}
		if len(t) > 0 { s.Items = infer(t[0]) }
		return s
	case string: return &schema{Type: "string"}
	case float64:
		if t == float64(int64(t)) { return &schema{Type: "integer"} }
		return &schema{Type: "number"}
	case bool: return &schema{Type: "boolean"}
	case nil: return &schema{Type: "null"}
	}
	return &schema{Type: "unknown"}
}

func mergeSchema(a, b *schema) *schema {
	if a == nil { return b }
	if b == nil { return a }
	if a.Type != b.Type { return &schema{Type: "any"} }
	if a.Type == "object" && a.Properties != nil && b.Properties != nil {
		merged := &schema{Type: "object", Properties: map[string]*schema{}}
		for k, v := range a.Properties { merged.Properties[k] = v }
		for k, v := range b.Properties {
			if existing, ok := merged.Properties[k]; ok { merged.Properties[k] = mergeSchema(existing, v) } else { merged.Properties[k] = v }
		}
		return merged
	}
	return a
}

func main() {
	var r io.Reader = os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		defer f.Close(); r = f
	}

	data, _ := io.ReadAll(r)
	content := strings.TrimSpace(string(data))

	var result *schema
	if strings.HasPrefix(content, "{") || strings.HasPrefix(content, "[") {
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		result = infer(v)
	} else {
		// JSONL
		for _, line := range strings.Split(content, "\n") {
			line = strings.TrimSpace(line); if line == "" { continue }
			var v interface{}
			if err := json.Unmarshal([]byte(line), &v); err != nil { continue }
			result = mergeSchema(result, infer(v))
		}
	}

	// Add $schema
	out := map[string]interface{}{"$schema": "http://json-schema.org/draft-07/schema#"}
	b, _ := json.Marshal(result)
	var s map[string]interface{}; json.Unmarshal(b, &s)
	for k, v := range s { out[k] = v }
	b, _ = json.MarshalIndent(out, "", "  "); fmt.Println(string(b))
}
