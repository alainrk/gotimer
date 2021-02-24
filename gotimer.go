package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/nwillc/gotimer/gen/version"
	"github.com/nwillc/gotimer/setup"
	"github.com/nwillc/gotimer/typeface"
	"github.com/nwillc/gotimer/utils"
	"os"
	"time"
)

func main() {
	flags := &setup.Values{}
	setup.NewFlagSetWithValues(os.Args[0], flags).ParseOsArgs()
	if *flags.Version {
		fmt.Println("Version:", version.Version)
		os.Exit(0)
	}

	duration, err := time.ParseDuration(*flags.Time)
	if err != nil {
		panic(err)
	}

	color := tcell.ColorNames[*flags.ColorName]
	var s tcell.Screen
	if s, err = tcell.NewScreen(); err != nil {
		panic(err)
	}

	if err := s.Init(); err != nil {
		panic(err)
	}

	// paused indicates timer is paused
	var paused = false
	go func() {
		for {
			time.Sleep(time.Second)
			if paused {
				continue
			}
			display(duration, s, color, *flags.FontName)
			duration = duration - time.Second
			if duration < 0 {
				_ = s.Beep()
				break
			}
		}
	}()

	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				os.Exit(0)
			} else if ev.Rune() == ' ' {
				paused = !paused
			}
		}
	}
}

func display(duration time.Duration, s tcell.Screen, color tcell.Color, fontName string) {
	font, ok := typeface.AvailableFonts[fontName]
	if !ok {
		panic("font not available")
	}
	s.Clear()
	str, err := utils.Format(duration)
	if err != nil {
		panic(err)
	}
	x := 1
	for _, c := range str {
		width, err := typeface.RenderRune(s, c, font, color, x, 1)
		if err != nil {
			panic(err)
		}
		x += width + 1
	}
	s.Show()
}
