// A simple package for adding colourful comments onto source code lines
// primarily for the use of user-friendly error messages.
package decorator

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// An enumeration defining how to style a line.
type LineColourEnum string

// The different styles that can be applied to a line.
// Currently you can only specify one value per line.
const (
	Normal    = "\033[0m"
	Bold      = "\033[1m"
	FgBlack   = "\033[30m"
	FgRed     = "\033[31m"
	FgGreen   = "\033[32m"
	FgYellow  = "\033[33m"
	FgBlue    = "\033[34m"
	FgMagenta = "\033[35m"
	FgCyan    = "\033[36m"
	FgWhite   = "\033[37m"
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

type commentInfo struct {
	at      int
	text    string
	colours []LineColour
}

type lineInfo struct {
	text           string
	colours        []LineColour
	meta           LineMetadata
	topComments    []commentInfo
	bottomComments []commentInfo
}

// Describes how to colour a particular segment of a line.
type LineColour struct {
	// The index of the first character to colour.
	From int

	// The index of the last character (non-inclusive) to colour.
	To int

	// The colouring to apply.
	Colour LineColourEnum
}

// Describes metadata about a line.
type LineMetadata struct {
	// The line number that this line was extracted from.
	LineNumber int

	// The file that this line belongs to.
	FileName     string
	cachedPrefix string
}

// The main type responsible for this library's functionality.
type Decorator struct {
	lines []lineInfo
}

// Adds a new line to be decorated.
func (dec *Decorator) AddLine(line string, meta LineMetadata) error {
	if strings.ContainsAny(line, "\n\t\r") {
		return errors.New("string contains one of ['\\n', '\\t', '\\r'] which aren't supported")
	}

	dec.lines = append(dec.lines, lineInfo{text: line, meta: meta})
	return nil
}

// Adds a comment below the specified line, pointing at a specific character in that line.
// Lines can have multiple bottom comments.
func (dec *Decorator) AddBottomComment(line int, at int, comment string) error {
	if strings.ContainsAny(comment, "\n\t\r") {
		return errors.New("string contains one of ['\\n', '\\t', '\\r'] which aren't supported")
	} else if line < 0 || line >= len(dec.lines) {
		return errors.New("index out of bounds")
	}

	dec.lines[line].bottomComments = append(dec.lines[line].bottomComments, commentInfo{text: comment, at: at})
	return nil
}

// Adds a comment above the specified line, pointing at a specific character in that line.
// Lines can have multiple top comments.
func (dec *Decorator) AddTopComment(line int, at int, comment string) error {
	if strings.ContainsAny(comment, "\n\t\r") {
		return errors.New("string contains one of ['\\n', '\\t', '\\r'] which aren't supported")
	} else if line < 0 || line >= len(dec.lines) {
		return errors.New("index out of bounds")
	}

	dec.lines[line].topComments = append(dec.lines[line].topComments, commentInfo{text: comment, at: at})
	return nil
}

// Applies colouring to the specified line.
func (dec *Decorator) ColourLine(line int, colour LineColour) error {
	if line < 0 || line >= len(dec.lines) {
		return errors.New("index out of bounds")
	}

	dec.lines[line].colours = append(dec.lines[line].colours, colour)
	return nil
}

// Applies colouring to the specified comment for the specified line.
func (dec *Decorator) ColourBottomComment(line int, comment int, colour LineColour) error {
	if line < 0 || line >= len(dec.lines) {
		return errors.New("line index out of bounds")
	} else if comment < 0 || comment >= len(dec.lines[line].bottomComments) {
		return errors.New("comment index out of bounds")
	}

	dec.lines[line].bottomComments[comment].colours = append(dec.lines[line].bottomComments[comment].colours, colour)
	return nil
}

// Applies colouring to the specified comment for the specified line.
func (dec *Decorator) ColourTopComment(line int, comment int, colour LineColour) error {
	if line < 0 || line >= len(dec.lines) {
		return errors.New("line index out of bounds")
	} else if comment < 0 || comment >= len(dec.lines[line].topComments) {
		return errors.New("comment index out of bounds")
	}

	dec.lines[line].topComments[comment].colours = append(dec.lines[line].topComments[comment].colours, colour)
	return nil
}

// Constructs a string consisting of all the lines; their metadata, and their comments.
func (dec *Decorator) String() string {
	var b strings.Builder

	// Generate prefixes and find the longest one
	longestPrefixLength := 0
	for i := 0; i < len(dec.lines); i++ {
		line := &dec.lines[i]
		line.meta.generatePrefix()
		if len(line.meta.cachedPrefix) > longestPrefixLength {
			longestPrefixLength = len(line.meta.cachedPrefix)
		}
	}

	// Write out each line and its comments + colours
	for i := 0; i < len(dec.lines); i++ {
		line := &dec.lines[i]

		// Write out top comments
		writeTopComments(&b, line, longestPrefixLength)

		// Write out line
		writePrefix(&b, line.meta.cachedPrefix, longestPrefixLength)
		writeColoured(&b, line.text, line.colours)
		b.WriteByte('\n')

		// Write out bottom comments
		writeBottomComments(&b, line, longestPrefixLength)
	}

	return b.String()
}

func writePrefix(b *strings.Builder, prefix string, longest int) {
	b.WriteString(prefix)
	writePadding(b, longest-len(prefix))
	b.WriteString(" | ")
}

func writePadding(b *strings.Builder, amount int) {
	for i := 0; i < amount; i++ {
		b.WriteByte(' ')
	}
}

func writeBottomComments(b *strings.Builder, line *lineInfo, longestPrefixLength int) {
	commentsWritten := 0
	for j := 0; j < len(line.bottomComments)*3; j++ {
		writePrefix(b, "", longestPrefixLength)

		cursor := 0
		mod := j % 3
		written := false

		for k := commentsWritten; k < len(line.bottomComments); k++ {
			comment := line.bottomComments[k]
			if cursor > comment.at {
				continue
			}
			writePadding(b, comment.at-cursor)
			cursor = comment.at

			if k == commentsWritten && !written {
				if mod == 0 {
					b.WriteRune('│')
					cursor++
				} else if mod == 1 {
					b.WriteByte('v')
					cursor++
				} else {
					written = true // Stop the other comments from acting like they need to be written out on this line.
					commentsWritten++
					writeColoured(b, comment.text, comment.colours)
					cursor += len(comment.text) // Stop the other comments from overwriting our text with their pipes.
				}
			} else {
				b.WriteRune('│')
				cursor++
			}
		}

		b.WriteByte('\n')
	}
}

func writeTopComments(b *strings.Builder, line *lineInfo, longestPrefixLength int) {
	commentsWritten := 0
	for j := 0; j < len(line.topComments)*3; j++ {
		writePrefix(b, "", longestPrefixLength)

		cursor := 0
		mod := j % 3
		written := false

		for k := 0; k < len(line.topComments); k++ {
			comment := line.topComments[k]
			if cursor > comment.at {
				continue
			}
			writePadding(b, comment.at-cursor)
			cursor = comment.at

			if k == commentsWritten && !written && mod == 0 {
				writeColoured(b, comment.text, comment.colours)
				written = true
				commentsWritten++
				cursor += len(comment.text)
			} else if k == commentsWritten-1 && mod == 1 {
				b.WriteByte('^')
				cursor++
			} else if k < commentsWritten {
				b.WriteRune('│')
				cursor++
			}
		}

		b.WriteByte('\n')
	}
}

func writeColoured(b *strings.Builder, text string, colours []LineColour) {
	sort.Slice(colours, func(i, j int) bool {
		return colours[i].From < colours[j].From
	})
	colourI := 0
	start := 0
	for colourI < len(colours) {
		colour := colours[colourI]
		colourI++

		if start > colour.From {
			colour.From = start // colour is a copy
		}
		b.WriteString(text[start:colour.From])
		start = colour.To

		b.WriteString(string(colour.Colour))
		b.WriteString(text[colour.From:colour.To])
		b.WriteString(Normal)
	}
	if start < len(text) {
		b.WriteString(text[start:])
	}
}

func (meta *LineMetadata) generatePrefix() {
	if meta.cachedPrefix != "" {
		return
	}
	meta.cachedPrefix = fmt.Sprintf("%s @ %d", meta.FileName, meta.LineNumber)
}
