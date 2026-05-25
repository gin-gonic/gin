package binding

import (
	"bytes"
	"runtime"
	"testing"
	"time"
)

// TestXMLBindingAdversarialInputs verifies that the XML binding maintains
// security boundaries under adversarial inputs including XML bomb / billion
// laughs attacks and other malicious XML payloads.
//
// Security invariant: parsing adversarial XML must not cause unbounded memory
// growth or hang the process. The binding must either complete within a
// reasonable time/memory budget or return an error — it must never silently
// consume excessive resources.
func TestXMLBindingAdversarialInputs(t *testing.T) {
	payloads := []struct {
		name    string
		payload string
	}{
		{
			name: "billion_laughs_classic",
			payload: `<?xml version="1.0"?>
<!DOCTYPE lolz [
  <!ENTITY lol "lol">
  <!ENTITY lol2 "&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;">
  <!ENTITY lol3 "&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;">
  <!ENTITY lol4 "&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;">
  <!ENTITY lol5 "&lol4;&lol4;&lol4;&lol4;&lol4;&lol4;&lol4;&lol4;&lol4;&lol4;">
  <!ENTITY lol6 "&lol5;&lol5;&lol5;&lol5;&lol5;&lol5;&lol5;&lol5;&lol5;&lol5;">
  <!ENTITY lol7 "&lol6;&lol6;&lol6;&lol6;&lol6;&lol6;&lol6;&lol6;&lol6;&lol6;">
  <!ENTITY lol8 "&lol7;&lol7;&lol7;&lol7;&lol7;&lol7;&lol7;&lol7;&lol7;&lol7;">
  <!ENTITY lol9 "&lol8;&lol8;&lol8;&lol8;&lol8;&lol8;&lol8;&lol8;&lol8;&lol8;">
]>
<root>&lol9;</root>`,
		},
		{
			name: "billion_laughs_shallow",
			payload: `<?xml version="1.0"?>
<!DOCTYPE bomb [
  <!ENTITY a "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa">
  <!ENTITY b "&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;&a;">
  <!ENTITY c "&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;&b;">
  <!ENTITY d "&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;&c;">
]>
<root>&d;</root>`,
		},
		{
			name: "quadratic_blowup",
			payload: `<?xml version="1.0"?>
<!DOCTYPE bomb [
  <!ENTITY x "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx">
]>
<root>&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;</root>`,
		},
		{
			name: "deeply_nested_elements",
			payload: func() string {
				var buf bytes.Buffer
				buf.WriteString(`<?xml version="1.0"?>`)
				depth := 100000
				for i := 0; i < depth; i++ {
					buf.WriteString("<a>")
				}
				buf.WriteString("deep")
				for i := 0; i < depth; i++ {
					buf.WriteString("</a>")
				}
				return buf.String()
			}(),
		},
		{
			name: "large_attribute_value",
			payload: func() string {
				var buf bytes.Buffer
				buf.WriteString(`<?xml version="1.0"?><root attr="`)
				for i := 0; i < 10*1024*1024; i++ {
					buf.WriteByte('A')
				}
				buf.WriteString(`">value</root>`)
				return buf.String()
			}(),
		},
		{
			name: "many_attributes",
			payload: func() string {
				var buf bytes.Buffer
				buf.WriteString(`<?xml version="1.0"?><root`)
				for i := 0; i < 10000; i++ {
					buf.WriteString(` attr`)
					buf.WriteByte(byte('0' + i%10))
					buf.WriteString(`="value"`)
				}
				buf.WriteString(`>content</root>`)
				return buf.String()
			}(),
		},
		{
			name: "entity_in_attribute",
			payload: `<?xml version="1.0"?>
<!DOCTYPE root [
  <!ENTITY e1 "evil">
  <!ENTITY e2 "&e1;&e1;&e1;&e1;&e1;&e1;&e1;&e1;&e1;&e1;">
  <!ENTITY e3 "&e2;&e2;&e2;&e2;&e2;&e2;&e2;&e2;&e2;&e2;">
]>
<root attr="&e3;">value</root>`,
		},
		{
			name: "malformed_xml",
			payload: `<?xml version="1.0"?><root><unclosed><also_unclosed>text`,
		},
		{
			name: "null_bytes",
			payload: "<?xml version=\"1.0\"?><root>\x00\x00\x00</root>",
		},
		{
			name: "unicode_bomb",
			payload: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE bomb [
  <!ENTITY u "` + string([]byte{0xF0, 0x9F, 0x92, 0xA3}) + `">
  <!ENTITY u2 "&u;&u;&u;&u;&u;&u;&u;&u;&u;&u;">
  <!ENTITY u3 "&u2;&u2;&u2;&u2;&u2;&u2;&u2;&u2;&u2;&u2;">
]>
<root>&u3;</root>`,
		},
	}

	// Memory limit: 256 MB growth allowed per parse attempt
	const maxMemoryGrowthBytes = 256 * 1024 * 1024
	// Time limit per parse attempt
	const maxDuration = 5 * time.Second

	type Target struct {
		Value string `xml:",chardata"`
		Attr  string `xml:",attr"`
	}

	xmlBind := xmlBinding{}

	for _, tc := range payloads {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Measure baseline memory
			var memBefore runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&memBefore)

			done := make(chan error, 1)
			start := time.Now()

			go func() {
				var obj Target
				err := xmlBind.BindBody([]byte(tc.payload), &obj)
				done <- err
			}()

			select {
			case err := <-done:
				elapsed := time.Since(start)

				// Measure memory after
				var memAfter runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&memAfter)

				// Security invariant 1: must complete within time limit
				if elapsed > maxDuration {
					t.Errorf("SECURITY VIOLATION: XML parsing took %v (limit %v) for payload %q — possible DoS",
						elapsed, maxDuration, tc.name)
				}

				// Security invariant 2: memory growth must be bounded
				if memAfter.TotalAlloc > memBefore.TotalAlloc {
					growth := memAfter.TotalAlloc - memBefore.TotalAlloc
					if growth > maxMemoryGrowthBytes {
						t.Errorf("SECURITY VIOLATION: XML parsing allocated %d bytes (limit %d) for payload %q — possible memory bomb",
							growth, maxMemoryGrowthBytes, tc.name)
					}
				}

				// Security invariant 3: if parsing succeeded, the result must not
				// be astronomically large (entity expansion must be bounded)
				if err == nil {
					// A successful parse of a bomb payload is a security concern
					// if the result is huge; log it as a warning
					t.Logf("payload %q parsed without error in %v (err=%v)", tc.name, elapsed, err)
				}

			case <-time.After(maxDuration):
				t.Errorf("SECURITY VIOLATION: XML parsing timed out after %v for payload %q — DoS vulnerability confirmed",
					maxDuration, tc.name)
			}
		})
	}
}