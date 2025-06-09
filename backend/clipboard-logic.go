package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
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

func saveClipboardHistoy(){
	f,err := os.OpenFile("/home/anup/.clipboardHistory", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil{
		panic(err)
	}
	defer f.Close()

	jsonBytes,err := json.Marshal(clipboardHistory)
	if err != nil{
		panic(err)
	}
	
	_, err = f.WriteString(string(jsonBytes))
	if err != nil {
		panic(err)
	}
}

func main(){
	data,err := os.ReadFile("/home/anup/.clipboardHistory")
	if err!=nil{
		panic(err)
	}
	
	var arr []string
	err = json.Unmarshal(data,&arr)
	if err != nil{
		panic(err)
	}

	clipboardHistory = append(clipboardHistory, arr...)

	go connection()
	defer saveClipboardHistoy()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	stop := make(chan bool)
	
	go func(){
		for{
			select{
				case <-stop:
					return
				default:
					content, err := getClipboard()
					if err == nil  && strings.TrimSpace(content) != "" {
						addText(content)
					}
					time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	<-c
	close(stop)
}