package main

import (
	"os/exec"
	"strings"
	"time"
)

func getClipboard() (string, error) {
	out, err := exec.Command("wl-paste").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

var clipboardHistory = make([]string, 0 ,20)

func addText(text string){
	if len(clipboardHistory) > 0 && clipboardHistory[len(clipboardHistory)-1] == text{
		return
	}

	if len(clipboardHistory) == 20{
		clipboardHistory = clipboardHistory[1:20]
		clipboardHistory = append(clipboardHistory,text)
		return
	}

	clipboardHistory = append(clipboardHistory, text)
}

func getClipboardHistory() []string{
	return clipboardHistory
}

func main(){
	go connection()
	for{
		content, err := getClipboard()
		
		if err == nil  && strings.TrimSpace(content) != "" {
			addText(content)
		}

		time.Sleep(500 * time.Millisecond)
	}
}