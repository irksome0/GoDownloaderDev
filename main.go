package main

import (
	"fmt"
	"image/color"
	"io"
	"os"

	// Downloading videos
	"github.com/kkdai/youtube/v2"

	// GUI
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	
	// Initializing a new GUI application
	application := app.New()

	// Setting up basic settings 1. width 2. height 3. basic downloading path 4. the application
	mainWindow(480, 250, "D:/GoDownloader", application)

	// Runner
	application.Run()
}

func mainWindow(windowWidth float32, windowHeight float32, path string, application fyne.App) {

	// Setting up the window
	myWindow := application.NewWindow("Downloader")
	myWindow.Resize(fyne.NewSize(windowWidth, windowHeight))

	// Initializing welcome text
	text := canvas.NewText("Welcome to downloader", color.White)
	text.Move(fyne.NewPos(windowWidth*0.3, windowHeight/10))

	// Initializing video url input field
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter ulr...")
	input.Resize(fyne.NewSize(windowWidth*0.96, windowHeight/6.25))
	input.Move(fyne.NewPos(windowWidth/108, windowHeight/4))

	// Download button
	button := widget.NewButton("Download", func() {
		fmt.Printf("Your url: %s\n", input.Text)
		download(input.Text, path)
	})
	button.Resize(fyne.NewSize(windowWidth/4, windowHeight/5))
	button.Move(fyne.NewPos(windowWidth*0.355, windowHeight/2.25))

	// Button to open a window to change the path of downloading audio
	locationButton := widget.NewButton("Change path", func() {
		locationWindow(windowWidth, windowHeight, path, application)
		myWindow.Close()
	})
	locationButton.Resize(fyne.NewSize(windowWidth/4, windowHeight/5))
	locationButton.Move(fyne.NewPos(windowWidth*0.355, windowHeight/1.5))

	// Setting up all the content in the window
	myWindow.SetContent(container.NewWithoutLayout(text, input, button, locationButton))
	myWindow.Show()
}

func locationWindow(windowWidth float32, windowHeight float32, path string, application fyne.App) {

	// Setting up the window
	myWindow := application.NewWindow("Create a folder")
	myWindow.Resize(fyne.NewSize(windowWidth, windowHeight))

	// Initializing the description text
	text := canvas.NewText("Enter location for the folder", color.White)
	text.Move(fyne.NewPos(windowWidth*0.3, windowHeight/10))

	// Initializing download path input field
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter path...")
	input.Resize(fyne.NewSize(windowWidth*0.96, windowHeight/6.25))
	input.Move(fyne.NewPos(windowWidth/108, windowHeight/4))

	// Button to save the path
	button := widget.NewButton("Use path", func() {
		path = input.Text
	})
	button.Resize(fyne.NewSize(windowWidth/4, windowHeight/5))
	button.Move(fyne.NewPos(windowWidth*0.355, windowHeight/2.25))

	// Button to return to the main window
	returnButton := widget.NewButton("Return", func() {
		mainWindow(windowWidth, windowHeight, path, application)
		myWindow.Close()
	})
	returnButton.Resize(fyne.NewSize(windowWidth/4, windowHeight/4))
	returnButton.Move(fyne.NewPos(windowWidth*0.355, windowHeight/1.5))

	// Setting up all the content
	myWindow.SetContent(container.NewWithoutLayout(text, input, button, returnButton))
	myWindow.Show()
}

func download(vID string, path string) {
	var videoID string
	
	// Validating a video URL and saving its ID
	if len(vID) < 43 {
		return
	}
	if len(vID) == 43 {
		videoID = vID[32:]
	} else {
		videoID = vID[32:43]
	}

	// Initializing client to download videos
	client := youtube.Client{}
	video, err := client.GetVideo(videoID)
	if err != nil {
		fmt.Println("Error with downloading the video")
		return
	}
	
	// Displaying video information if it is valid in console
	// not necessary function
	displayInfo(*video)

	// Getting a video
	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}

	// Creating the path folder
	os.Mkdir(path, 0700)

	// Saving the audio
	file, err := os.Create(path + "/" + video.Title + ".mp4")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
}

func displayInfo(video youtube.Video) {
	fmt.Printf("Title: %s\n", video.Title)
	fmt.Printf("Author: %s\n", video.Author)
	fmt.Printf("Duration: %s\n", video.Duration.String())
}
