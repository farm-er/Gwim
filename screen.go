package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

const (
	normalStyle           = iota
	stringStyle           // "
	multiLinesStringStyle // `
	littleString          // '
)

var (
	style                       = normalStyle
	screenStyle     tcell.Style = tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorWhite)
	statusBarStyle  tcell.Style = tcell.StyleDefault.Bold(true).Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	lineNumberStyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorWhite)
	screen          tcell.Screen
)

func initScreen() {

	var err error
	screen, err = tcell.NewScreen()

	if err != nil {
		log.Fatalf("Error creating screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Error initializing screen: %v", err)
	}

	screen.SetStyle(screenStyle)

}

func refresh(c Cursor) {
	width, height = screen.Size()
	screen.Clear()

	view = cur.Line - (height - 2)
	if view < 0 {
		view = 0
	}
	cur.Start = len(strconv.Itoa(cur.Line+1)) + 3
	for i, contentLine := range c.Content[view:] {
		// craft the number of the line
		lineNumber := strconv.Itoa(i+view+1) + " ~ "
		lineStart := len(lineNumber)
		// draw the number of the line
		for l, digit := range lineNumber {
			screen.SetContent(l, i, rune(digit), nil, lineNumberStyle)
		}

		for j, ch := range contentLine {
			switch ch {
			case '"':
				if style == normalStyle {
					screenStyle = tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorDarkGreen.TrueColor())
					style = stringStyle
					screen.SetContent(j+lineStart, i, ch, nil, screenStyle)
				} else if style == stringStyle {
					screen.SetContent(j+lineStart, i, ch, nil, screenStyle)
					screenStyle = tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorWhite.TrueColor())
					style = normalStyle
				} else {
					screen.SetContent(j+lineStart, i, ch, nil, screenStyle)
				}
			case '`':
				if style == normalStyle {
					screenStyle = tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorIndianRed.TrueColor())
					style = multiLinesStringStyle
					screen.SetContent(j+lineStart, i, ch, nil, screenStyle)
				} else if style == multiLinesStringStyle {
					screen.SetContent(j+lineStart, i, ch, nil, screenStyle)
					screenStyle = tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorWhite.TrueColor())
					style = normalStyle
				} else {
					screen.SetContent(j+lineStart, i, ch, nil, screenStyle)
				}
			case '\'':
				if style == normalStyle {
					screenStyle = tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorDarkGreen.TrueColor())
					style = stringStyle
					screen.SetContent(j+lineStart, i, ch, nil, screenStyle)
				} else if style == stringStyle {
					screen.SetContent(j+lineStart, i, ch, nil, screenStyle)
					screenStyle = tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorWhite.TrueColor())
					style = normalStyle
				} else {
					screen.SetContent(j+lineStart, i, ch, nil, screenStyle)
				}
			case '{', '}', '[', ']', '*', '+', '=', '-', '%', '!', '<', '>', ':', '.', ',', '\\':
				screen.SetContent(j+lineStart, i, ch, nil, tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorDarkCyan.TrueColor()))
			case '|', '&', '(', ')':
				screen.SetContent(j+lineStart, i, ch, nil, tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorRebeccaPurple.TrueColor()))
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				screen.SetContent(j+lineStart, i, ch, nil, tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorLightBlue.TrueColor()))
			default:
				screen.SetContent(j+lineStart, i, ch, nil, screenStyle)
			}

		}
		// NEED TO OPTIMIZE THIS: WE DON'T NEED TO CHECK EVERY TIME THE POS AND LINE WE CAN ADD IT AFTER THE LOOP
		if i == cur.Line && cur.Pos < lineStart {
			cur.Pos = lineStart
		}
		if style != multiLinesStringStyle {
			screenStyle = tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorWhite)
			style = normalStyle
		}
	}
	screenStyle = tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorWhite)
	style = normalStyle
	status := fmt.Sprintf("lines: %v line: %v line-starts-at: %v length: %v position: %v screen: (w: %v,h: %v) ", len(c.Content), cur.Line+1, cur.Start, len(c.Content[cur.Line]), cur.Pos, width, height)
	statusRunes := []rune(status)
	for j, ch := range statusRunes {
		screen.SetContent(j, height-1, ch, nil, statusBarStyle)
	}
	screen.ShowCursor(cur.Pos, cur.Line-view)
	screen.Show()
}
