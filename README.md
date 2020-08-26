# polybar-spotify-go
Small script to get Spotify control to work with Polybar. This is almost identical to [the Python version](https://github.com/Sudneo/polybar-spotify), rewritten in GO and without the necessity of getting the song title every second.

Compared to the Python version this is 4-6 times faster and also uses Dbus signals to track song and playbackstatus changes, without polling.

## Spotify Control and Polybar

In order to have this to work you need to simply download a release or build it in your machine.

## Idea

The idea is to have a simple way to see Spotify music playing and simple controls to pause/play songs and to go to previous/next song.

## Limitations

The main limitation is that Spotify application doesn't expose the current position in the song,
so we cannot show a fancy progress bar of the song.
It could be possible to use a timer and signals to check when a new song started, but without information about changing the position in the song it would have several issues.

## Configuration

I deliberately wrote a script that could do all I needed and then split the functionalities in different polybar modules.

The result was something like the following:
```
[...]

modules-left = spotify-song spotify-backward spotify-status spotify-forward

[...]
[module/spotify-song]
type = custom/script
exec = ~/polybar-scripts/polybar-spotify-go -album
tail = true
format = <label>
format-foreground = #fff
format-background = #773f3f3f
format-underline = #c9665e
format-padding = 4
label = %{T4}%output%%{T-}

[module/spotify-backward]
type = custom/script
exec = ~/polybar-scripts/polybar-spotify-go -prevIcon
click-left = ~/polybar-scripts/polybar-spotify-go -prev
interval = 60
format-padding = -1

[module/spotify-status]
type = custom/script
exec = ~/polybar-scripts/polybar-spotify-go -playpause-icon
click-left = ~/polybar-scripts/polybar-spotify-go -playpause
tail = true
format-padding = -1

[module/spotify-forward]
type = custom/script
exec = ~/polybar-scripts/polybar-spotify-go -nextIcon
click-left = ~/polybar-scripts/polybar-spotify-go -next
interval = 60
format-padding = -1
```

Few things to notice:

* The custom font usage in the `spotify-song` block is not casual. Since the blocks are separated,
  the best result is obtained with monospace fonts, so that the width of the block is fixed.
* When Spotify is not Running, the script doesn't print anything at all, so the whole block will
  basically disappear.
* Depending on Polybar polling time, the "next" and "previous" icons might take a while to come up when Spotify gets launched once Polybar is running.

## Outcome

When a song is playing, the result is like this:

![Playing](playing.png)

When a song is paused, the result is like this:

![Paused](paused.png)

Note that the script trims or pads the output to a fixed amount of characters, that together with a monospace font produce the result of a fixed width polybar block.

