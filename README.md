# EbitUI

A simple UI library written in [Go](https://go.dev/) for [Ebitengine](https://ebitengine.org/).

## Development is just begining. Don't use this unless you are insane.
## TO DO: Pretty much the whole thing.



![Screenshot of some basic code and a basic window](Example/example.png)

Example is in the 'Example' directory.

## EbitUI.Start()

Init EbitUI


## Window IDs (string)

These can be whatever you like, but they are converted to lowercase.

## EbitUI.DrawWindows(screen)

Render windows to the screen. Use this in Draw() at the end.


## EbitUI.AddWindow("window id", windowData)

Adds a window called "window id" to the window list


## EbitUI.DeleteWindow("window id")

Deletes a window called "window id" from the window list (if found)


## EbitUI.OpenWindow("window id")

Opens the window "window id" (if found)


## EbitUI.CloseWindow("window id")

Closes the window "window id" (if found)


## EbitUI.UpdateViewerSize(width, height int)

Updates the viewer size. This will prevent windows from going off the screen and will resize the hud.