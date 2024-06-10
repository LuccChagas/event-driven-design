package main

import (
	"fmt"
	"sync"
)

type ResizeEvent struct {
	Width  int
	Height int
}

type ResizeEventHandler func(event ResizeEvent) (int, int, error)

// Creation and Instance of Listeners
// They are responsable for "watch/wait" events
type ResizeEventListener struct {
	Events  chan ResizeEvent
	Handler ResizeEventHandler
	wg      *sync.WaitGroup
}

func NewResizeEventListener(handler ResizeEventHandler, wg *sync.WaitGroup) ResizeEventListener {
	return ResizeEventListener{make(chan ResizeEvent), handler, wg}
}

// Event Types
func (l *ResizeEventListener) Start(w *Window) {
	go func() {
		for event := range l.Events {
			new_width, new_height, err := l.Handler(event)
			if err != nil {
				break
			}
			w.width = new_width
			w.Height = new_height
			l.wg.Done()
		}

	}()

}

func (l *ResizeEventListener) Stop() {
	close(l.Events)
}

func (l *ResizeEventListener) Send(event ResizeEvent) {
	l.wg.Add(1)
	l.Events <- event
}

type Window struct {
	listener ResizeEventListener
	title    string
	width    int
	Height   int
}

// Now I'll Create the object that could be manipulated (the object who can execute the events) In our case the Window
func CreateWindow(title string, width int, height int, handler ResizeEventHandler, wg *sync.WaitGroup) Window {
	return Window{listener: NewResizeEventListener(handler, wg), title: title, width: width, Height: height}
}

func (w *Window) Open() {
	w.listener.Start(w)
	fmt.Printf("Window %s opened with size %dx%d\n", w.title, w.width, w.Height)
}

func (w *Window) Close() {
	w.listener.Stop()
	fmt.Printf("Window %s closed\n", w.title)
}

func (w *Window) Resize(event ResizeEvent) {
	w.listener.Send(event)
}

func main() {
	var wg sync.WaitGroup

	window := CreateWindow("My Window", 800, 600, func(event ResizeEvent) (int, int, error) {
		fmt.Printf("Window resized to %dx%d\n", event.Width, event.Height)
		return event.Width, event.Height, nil
	}, &wg)

	window.Open()
	window.Resize(ResizeEvent{1024, 768})
	window.Resize(ResizeEvent{644, 494})
	window.Resize(ResizeEvent{800, 600})
	window.Resize(ResizeEvent{200, 120})
	window.Resize(ResizeEvent{900, 1144})
	window.Resize(ResizeEvent{500, 100})
	window.Resize(ResizeEvent{2140, 4116})

	wg.Wait()
	window.Close()

	fmt.Printf("Height is %d\n", window.Height)
	fmt.Printf("width is %d\n", window.width)
}
