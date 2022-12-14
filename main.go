package main

import (
    "fmt"
    "runtime"
    "time"

    "github.com/gotk3/gotk3/glib"
    "github.com/gotk3/gotk3/gtk"
)

func animatedText(frame int) string {
    var sym rune

    switch frame % 4 {
    case 0: sym = '-'
    case 1: sym = '/'
    case 2: sym = '-'
    case 3: sym = '\\'
    }

    text :=            " Очень сложная анимация \n\n\n"

    text += fmt.Sprintf("                  %c\n\n\n", sym)

    return text
}

func makeAnimatedWindow(stopped chan<- struct{}, button_enabled chan<- bool) {
    button_enabled <- false

    win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

    stop := make(chan struct{}, 1)
    win.Connect("destroy", func() {
        stop <- struct{}{}
    })

    label, _ := gtk.LabelNew(animatedText(0))

    win.Add(label)

    win.ShowAll()

    go func() {
        frame := 1

        for {
            time_chan := time.After(time.Second)

            select {
            case <-stop:
                select {
                case <-time_chan:
                    stopped <- struct{}{}
                    button_enabled <- true
                    return
                }

            case <-time_chan:
            }

            currentFrame := frame

            glib.IdleAdd(func() {
                label.SetLabel(animatedText(currentFrame))
            })

            frame += 1
        }
    }()
}


func main() {
    runtime.LockOSThread()

    gtk.Init(nil)

    win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    win.Connect("destroy", func() {
        gtk.MainQuit()
    })

    win.SetDefaultSize(400, 400)

    button, _ := gtk.ButtonNewWithLabel("Show animated window")
    win.Add(button)
    win.ShowAll()

    stopped := make(chan struct{}, 1)
    stopped <- struct{}{}

    button_enabled := make(chan bool, 1)

    button.Connect("clicked", func() {
        select {
        case <-stopped: makeAnimatedWindow(stopped, button_enabled)
        default:
        }
    })

    go func(){
        for {
            select {
            case state := <-button_enabled:
                glib.IdleAdd(func() {
                    button.SetSensitive(state)
                })
            }
        }
    }()

    gtk.Main()
}
