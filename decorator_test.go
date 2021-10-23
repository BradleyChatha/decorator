package decorator

import "testing"

func TestManual(t *testing.T) {
	var d Decorator

	d.AddLine("0123456789", LineMetadata{FileName: "some file", LineNumber: 69})
	d.ColourLine(0, LineColour{From: 2, To: 5, Colour: FgMagenta})
	d.ColourLine(0, LineColour{From: 4, To: 8, Colour: BgMagenta})

	d.AddBottomComment(0, 0, "abc")
	d.AddBottomComment(0, 5, "123lolol")
	d.AddBottomComment(0, 9, "doe ray me")
	d.ColourBottomComment(0, 0, LineColour{From: 0, To: 3, Colour: BgCyan})

	d.AddTopComment(0, 2, "abc")
	d.AddTopComment(0, 7, "123lolol")
	d.AddTopComment(0, 9, "doe ray me")
	d.ColourTopComment(0, 1, LineColour{From: 0, To: 3, Colour: FgGreen})
	d.ColourTopComment(0, 1, LineColour{From: 3, To: 6, Colour: FgRed})

	t.Log("\n\n\b", d.String())
}
