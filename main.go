package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
)

func trim_or_pad(s string, n int) string {
	if len(s) > n {
		return s[:n]
	} else {
		pad := ""
		for i := 0; i < (n - len(s)); i++ {
			pad = pad + " "
		}
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
	switch {
	case *next:
		obj.Call("org.mpris.MediaPlayer2.Player.Next", 0)
		os.Exit(0)
	case *nextIcon:
		fmt.Printf("\uf051\n")
		os.Exit(0)
	case *prev:
		obj.Call("org.mpris.MediaPlayer2.Player.Previous", 0)
		os.Exit(0)
	case *prevIcon:
		fmt.Printf("\uf048\n")
		os.Exit(0)
	case *playPause:
		obj.Call("org.mpris.MediaPlayer2.Player.PlayPause", 0)
		os.Exit(0)
	default:
		currentStatus, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get property:", err)
			os.Exit(1)
		}
		status := currentStatus.String()
		if status != "\"Playing\"" && status != "\"Paused\"" && status != "\"Stopped\"" {
			fmt.Println("Spotify is not running")
			os.Exit(0)
		}
		if *playPauseIcon {
			if status != "\"Playing\"" {
				fmt.Printf("\uf04b\n")
			} else {
				fmt.Printf("\uf04c\n")
			}
			os.Exit(0)
		}
		// Here we are subscribing to dbus signal and printing the new songs as they come
		metadata, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
		if err != nil {
			fmt.Println(err)
		} else {
			values := metadata.Value()
			title := values.(map[string]dbus.Variant)["xesam:title"]
			artist := values.(map[string]dbus.Variant)["xesam:artist"].Value().([]string)[0]
			album := values.(map[string]dbus.Variant)["xesam:album"]
			dash := "\u2014"
			status_icon := "\uf1bc"
			status_string := fmt.Sprintf("%s %s %s %s (%s)", status_icon, title, dash, artist, album)
			fmt.Println(trim_or_pad(status_string, *justify))
		}
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
			} else {
				values := metadata.Value()
				title := values.(map[string]dbus.Variant)["xesam:title"]
				artist := values.(map[string]dbus.Variant)["xesam:artist"].Value().([]string)[0]
				album := values.(map[string]dbus.Variant)["xesam:album"]
				dash := "\u2014"
				status_icon := "\uf1bc"
				status_string := fmt.Sprintf("%s %s %s %s (%s)", status_icon, title, dash, artist, album)
				fmt.Println(trim_or_pad(status_string, *justify))
			}
		}
	}
}
