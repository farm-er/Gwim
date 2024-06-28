package main

import (
	"slices"
	"strconv"
	"unicode/utf8"
)

type Cursor struct {
	Pos     int
	Line    int
	Start   int
	TabSize int
	Content [][]rune
}

func (c *Cursor) moveUp() {
	if c.Line > 0 {
		c.Line--
		c.Start = len(strconv.Itoa(c.Line+1)) + 3
		if len(c.Content[c.Line]) < c.Pos {
			c.Pos = c.Start + len(c.Content[c.Line])
		}
	}
}

func (c *Cursor) moveDown() {
	if c.Line < len(c.Content)-1 {
		c.Line++
		c.Start = len(strconv.Itoa(c.Line+1)) + 3
		if len(c.Content[c.Line]) < c.Pos-c.Start {
			c.Pos = len(c.Content[c.Line]) + c.Start
		}
	}
}

func (c *Cursor) moveLeft() {
	if c.Pos > c.Start {
		c.Pos--
	} else if c.Pos == c.Start {
		if c.Line > 0 {
			c.Line--
			c.Pos = len(c.Content[c.Line]) + c.Start
		}
	}
}

func (c *Cursor) moveRight() {
	if c.Pos-(c.Start) < len(c.Content[c.Line]) {
		c.Pos++
	}
}

func (c *Cursor) DoubleWrite(chars []rune) {
	if c.Pos == c.Start {
		c.Content[c.Line] = append(chars, c.Content[c.Line]...)
		c.Pos++
	} else if c.Pos == len(c.Content[c.Line])+c.Start {
		c.Content[c.Line] = append(c.Content[c.Line], chars...)
		c.Pos++
	} else if c.Pos > c.Start && c.Pos < len(c.Content[c.Line])+c.Start {
		c.Content[c.Line] = slices.Insert(c.Content[c.Line], c.Pos-c.Start, chars...)
		c.Pos++
	}
}

func (c *Cursor) Tab() {
	for i := 0; i < c.TabSize; i++ {
		c.Content[c.Line] = append(c.Content[c.Line], rune(' '))
	}
	c.Pos += c.TabSize
}

func (c *Cursor) Delete() {
	if c.Pos == c.Start {
		if c.Line > 0 {
			// we delete the empty line you were in
			// update the line position to put you in the the previous line
			// delete '\n' character
			// update your position
			lineToMove := c.Content[c.Line]
			c.Content = slices.Delete(c.Content, c.Line, c.Line+1)

			c.Line--
			c.Content[c.Line] = append(c.Content[c.Line], lineToMove...)
			c.Pos = len(strconv.Itoa(c.Line+1)) + 3 + len(c.Content[c.Line])
		}
	} else if c.Pos > c.Start && c.Pos < len(c.Content[c.Line])+c.Start+1 {
		// we delete the content you wanted to delete
		// and update your position
		c.Content[c.Line] = slices.Delete(c.Content[c.Line], c.Pos-c.Start-1, c.Pos-c.Start)
		c.Pos--
	}
}

func (c *Cursor) Enter() {
	if c.Pos == c.Start {
		newLine := c.Content[c.Line]
		c.Content = append(c.Content, []rune{})
		c.Content[c.Line] = []rune{}
		c.Line++
		c.Content[c.Line] = newLine
		c.Pos = len(strconv.Itoa(c.Line+1)) + 3 + len(c.Content[c.Line])
	} else if c.Pos == len(c.Content[c.Line])+c.Start {
		c.Content = slices.Insert(c.Content, c.Line+1, []rune{})
		c.Line++
		c.Pos = len(strconv.Itoa(c.Line+1)) + 3
	} else if c.Pos > c.Start && c.Pos < len(c.Content[c.Line])+c.Start {
		newLine := append([]rune{}, c.Content[c.Line][c.Pos-c.Start:]...)
		c.Content = slices.Insert(c.Content, c.Line+1, []rune(newLine))
		c.Content[c.Line] = c.Content[c.Line][:c.Pos-c.Start]
		c.Line++
		c.Content[c.Line] = newLine
		c.Pos = len(strconv.Itoa(c.Line+1)) + 3
	}
}

func (c *Cursor) AddChar(char rune) {
	if c.Pos == c.Start {
		newLine := []rune{char}
		newLine = append(newLine, c.Content[c.Line]...)
		c.Content[c.Line] = newLine
		c.Pos++
	} else if c.Pos == len(c.Content[c.Line])+c.Start {
		c.Content[c.Line] = append(c.Content[c.Line], char)
		c.Pos++
	} else if c.Pos > c.Start && c.Pos < len(c.Content[c.Line])+c.Start {
		c.Content[c.Line] = slices.Insert(c.Content[c.Line], c.Pos-c.Start, char)
		c.Pos++
	}
}

func (c *Cursor) GetFileContent(fileContent []byte) {
	i := 0
	for len(fileContent) > 0 {
		char, size := utf8.DecodeRune(fileContent)
		fileContent = fileContent[size:]
		if char == '\n' {
			c.Content = append(c.Content, []rune{})
			i++
		} else {
			c.Content[i] = append(c.Content[i], char)
		}
	}
}
