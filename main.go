package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/farm-er/text-editor-1/logging"
	"github.com/gdamore/tcell/v2"
)

var (
	cur    Cursor
	quit   bool = false
	view   int
	height int
	width  int
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("You need to provide a file path")
		os.Exit(0)
	}

	if len(os.Args) > 2 {
		fmt.Println("Too many arguments")
		os.Exit(0)
	}
	// NEED TO PROCESS THE FILE NAME FOR MORE CONTROL

	fileContent, err := os.ReadFile(os.Args[1])

	if err != nil {
		file, err := os.Create(os.Args[1])
		if err != nil {
			fmt.Println("Error creating the file")
			os.Exit(0)
		}
		defer file.Close()
		fileContent, err = io.ReadAll(file)

		if err != nil {
			fmt.Println("Error fetching file content")
		}
	}
	cur = Cursor{
		Line:    0,
		Pos:     0,
		Start:   0,
		TabSize: 4,
		Content: make([][]rune, 1),
	}

	cur.GetFileContent(fileContent)

	initScreen()
	logging.InitLogging()

	for !quit {
		refresh(cur)
		event := screen.PollEvent()
		switch ev := event.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				quit = true
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				cur.Delete()
			case tcell.KeyEnter:
				cur.Enter()
			case tcell.KeyTab:
				cur.Tab()
			case tcell.KeyUp:
				cur.moveUp()
			case tcell.KeyDown:
				cur.moveDown()
			case tcell.KeyLeft:
				cur.moveLeft()
			case tcell.KeyRight:
				cur.moveRight()
			case tcell.KeyRune:
				char := ev.Rune()

				switch char {
				case '(':
					cur.DoubleWrite([]rune{'(', ')'})
				case '{':
					cur.DoubleWrite([]rune{'{', '}'})
				case '[':
					cur.DoubleWrite([]rune{'[', ']'})
				case '&':
					cur.DoubleWrite([]rune{'&', '&'})
					cur.Pos++
				case '|':
					cur.DoubleWrite([]rune{'|', '|'})
					cur.Pos++
				case '"':
					cur.DoubleWrite([]rune{'"', '"'})
				case '\'':
					cur.DoubleWrite([]rune{'\'', '\''})
				case '`':
					cur.DoubleWrite([]rune{'`', '`'})
				default:
					cur.AddChar(char)
				}
			case tcell.KeyCtrlC:
				// COPY THE CONTENT
			case tcell.KeyCtrlS:
				// SAVE THE CONTENT
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}

	screen.Clear()
	screen.Fini()
	log.Println("Exiting gwim")
}
