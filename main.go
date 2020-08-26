package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
)

const (
	next         = "\uf051"
	prev         = "\uf048"
	pause        = "\uf04b"
	play         = "\uf04c"
	dash         = "\u2014"
	spotify_icon = "\uf1bc"
)

func trim_or_pad(s string, n int) string {
	fmt.Printf("Printing string of %v chars to %v\n", len(s), n)
	if len(s) > n {
		return s[:n]
	} else {
		pad := ""
		for i := 0; i < (n - len(s)); i++ {
			pad = pad + " "
		}
		fmt.Printf("Length of padded: %v\n", len(s+pad))
		return s + pad

	}
}

func main() {
	var playPause = flag.Bool("playpause", false, "Toggle Play/Pause, depending on current status")
	var playPauseIcon = flag.Bool("playpause-icon", false, "Print the icon for play/pause")
	var next = flag.Bool("next", false, "Go to next song")
	var nextIcon = flag.Bool("nextIcon", false, "Print the next icon")
	var prev = flag.Bool("prev", false, "Go to previous song")
	var prevIcon = flag.Bool("prevIcon", false, "Print the prev icon")
	var justify = flag.Int("justify", 75, "Justifies the output to the specified number of characters, padding or trimming")
	flag.Parse()
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()
	obj := conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2")
	currentStatus, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	if err != nil {
		os.Exit(1)
	}
	status := currentStatus.String()
	if status != "\"Playing\"" && status != "\"Paused\"" && status != "\"Stopped\"" {
		os.Exit(0)
	}
	switch {
	case *next:
		obj.Call("org.mpris.MediaPlayer2.Player.Next", 0)
		os.Exit(0)
	case *nextIcon:
		fmt.Println(next)
		os.Exit(0)
	case *prev:
		obj.Call("org.mpris.MediaPlayer2.Player.Previous", 0)
		os.Exit(0)
	case *prevIcon:
		fmt.Println(prev)
		os.Exit(0)
	case *playPause:
		obj.Call("org.mpris.MediaPlayer2.Player.PlayPause", 0)
		os.Exit(0)
	case *playPauseIcon:
		if status != "\"Playing\"" {
			fmt.Println(pause)
		} else {
			fmt.Println(play)
		}
		if err = conn.AddMatchSignal(
			dbus.WithMatchObjectPath("/org/mpris/MediaPlayer2"),
		); err != nil {
			panic(err)
		}
		c := make(chan *dbus.Signal, 10)
		conn.Signal(c)
		for _ = range c {
			currentStatus, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to get property:", err)
				os.Exit(1)
			}
			status := currentStatus.String()
			if status != "\"Playing\"" {
				fmt.Println(pause)
			} else {
				fmt.Println(play)
			}
		}
	default:
		metadata, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		values := metadata.Value()
		title := values.(map[string]dbus.Variant)["xesam:title"]
		artist := values.(map[string]dbus.Variant)["xesam:artist"].Value().([]string)[0]
		album := values.(map[string]dbus.Variant)["xesam:album"]
		status_string := fmt.Sprintf("%s %s %s %s (%s)", spotify_icon, title, dash, artist, album)
		fmt.Println(trim_or_pad(status_string, *justify))
		if err = conn.AddMatchSignal(
			dbus.WithMatchObjectPath("/org/mpris/MediaPlayer2"),
		); err != nil {
			panic(err)
		}
		c := make(chan *dbus.Signal, 10)
		conn.Signal(c)
		for _ = range c {
			metadata, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			values := metadata.Value()
			title := values.(map[string]dbus.Variant)["xesam:title"]
			artist := values.(map[string]dbus.Variant)["xesam:artist"].Value().([]string)[0]
			album := values.(map[string]dbus.Variant)["xesam:album"]
			status_string := fmt.Sprintf("%s %s %s %s (%s)", spotify_icon, title, dash, artist, album)
			fmt.Println(trim_or_pad(status_string, *justify))
		}
	}
}
