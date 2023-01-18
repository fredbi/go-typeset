package ansi

import "testing"

func TestStripToken(t *testing.T) {
	const token = "Lorem ipsum dolor \x1b[31;1;4msit amet,"
	/*
		remainder := "\x1b[31;1;4msit amet,"


		rstripped := StripANSIFromRunes([]rune(remainder))
		t.Logf("TEXT: %q", string(rstripped.Text))
		t.Logf("START: %q", string(rstripped.StartSequence))
		t.Logf("STOP: %q", string(rstripped.StopSequence))
		t.Logf("REM: %q", string(rstripped.Remainder))
	*/
	t.Logf("%v", StripToken([]rune(token)))
	/*
			    renderer_test.go:9: "Lorem ipsum dolor "
		    renderer_test.go:10: ""
		    renderer_test.go:11: ""
		    renderer_test.go:12: "\x1b[31;1;4msit amet,"
	*/
}
