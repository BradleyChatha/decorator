# Overview

Decorator is a simple library for adding comments onto lines of text, primarily aimed at user-friendly error messages.

For example:

![image](https://user-images.githubusercontent.com/3835574/138582615-02f3edb5-51ba-4af3-9edf-03fc901d2d8e.png)

The above output was generated from the following code:

```go
  var d Decorator
	
  d.AddLine("printf(\"%s\", 400)", LineMetadata{FileName: "main.c", LineNumber: 5})
  d.ColourLine(0, LineColour{From: 0, To: 6, Colour: FgYellow})
  d.ColourLine(0, LineColour{From: 7, To: 11, Colour: FgMagenta})
  d.ColourLine(0, LineColour{From: 13, To: 16, Colour: FgCyan})

  d.AddTopComment(0, 8, "%s was specified.")
  d.ColourTopComment(0, 0, LineColour{From: 0, To: 3, Colour: FgGreen})

  d.AddBottomComment(0, 13, "But a %d value was passed")
  d.ColourBottomComment(0, 0, LineColour{From: 6, To: 8, Colour: FgRed})
  
  print(d.String())
```

yada yada TODO
