package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func getClipboard() (string, error) {
	out, err := exec.Command("xclip", "-selection", "clipboard", "-o").Output()
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
	
	for i, str := range clipboardHistory {
		if str == text {
			clipboardHistory = append(clipboardHistory[:i], clipboardHistory[i+1:]...)
			break
		}
	}

	if len(clipboardHistory) == 20{
		clipboardHistory = clipboardHistory[1:]
	}

	clipboardHistory = append(clipboardHistory, text)
	log.Println("Text copied")
}

func getClipboardHistory() []string{
	return clipboardHistory
}

func saveClipboardHistoy(){
	log.Println("saving clipboard history...")
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
	log.Println("Daemon process starting...")
	log.Println("Reading clipboard history...")

	data, err := os.ReadFile("/home/anup/.clipboardHistory")

	log.Println("Reading clipboard history complete.")

	if err != nil{
		log.Panicln(err)
		panic(err)
	}

	var arr []string
	err = json.Unmarshal(data,&arr)
	
	if err != nil{
		arr = []string{}
	}

	clipboardHistory = append(clipboardHistory, arr...)

	go connection()
	defer saveClipboardHistoy()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	stop := make(chan bool)

	go func(){
		log.Println("Daemon started.")
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
	log.Panicln("Stopping daemon process.")
	close(stop)
}