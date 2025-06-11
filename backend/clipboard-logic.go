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

	data,err := os.ReadFile("/home/anup/.clipboardHistory")

	log.Println("Reading clipboard history complete.")

	if err!=nil{
		log.Panicln(err)
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
		log.Println("Daemon started.")
		for{
			select{
				case <-stop:
					return
				default:
					log.Println("Checking for new copy...")
					content, err := getClipboard()
					if err == nil  && strings.TrimSpace(content) != "" {
						log.Println("Copying new content to clipboard...")
						addText(content)
					}else{
						log.Println(err)
						log.Println("Not found.")
					}
					time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	<-c
	log.Panicln("Stopping daemon process.")
	close(stop)
}

// package main

// import (
// 	"bufio"
// 	"encoding/json"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"os/signal"
// 	"strings"
// 	"syscall"
// 	"time"
// )

// var clipboardHistory = make([]string, 0, 20)

// func addText(text string) {
// 	trimmedText := strings.TrimSpace(text)

// 	if trimmedText == "" || (len(clipboardHistory) > 0 && clipboardHistory[len(clipboardHistory)-1] == trimmedText) {
// 		return
// 	}

// 	if len(clipboardHistory) == 20 {
// 		clipboardHistory = clipboardHistory[1:20]
// 	}

// 	clipboardHistory = append(clipboardHistory, trimmedText)
// 	log.Printf("Text copied: %s (History size: %d)", trimmedText, len(clipboardHistory))
// }

// func getClipboardHistory() []string {
// 	return clipboardHistory
// }

// func saveClipboardHistory() {
// 	log.Println("Saving clipboard history...")
// 	f, err := os.OpenFile("/home/anup/.clipboardHistory", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
// 	if err != nil {
// 		log.Printf("ERROR: Could not open clipboard history file for saving: %v", err)
// 		return
// 	}
// 	defer f.Close()

// 	jsonBytes, err := json.Marshal(clipboardHistory)
// 	if err != nil {
// 		log.Printf("ERROR: Could not marshal clipboard history to JSON: %v", err)
// 		return
// 	}

// 	_, err = f.WriteString(string(jsonBytes))
// 	if err != nil {
// 		log.Printf("ERROR: Could not write clipboard history to file: %v", err)
// 		return
// 	}
// 	log.Println("Clipboard history saved successfully.")
// }

// func main() {
// 	log.Println("Daemon process starting...")
// 	log.Println("Reading clipboard history...")

// 	data, err := os.ReadFile("/home/anup/.clipboardHistory")

// 	if err!=nil{
// 		panic(err)
// 	}

// 	var loadedHistory []string
// 	err = json.Unmarshal(data, &loadedHistory)
// 	if err != nil {
// 		panic(err)
// 	} 
// 	clipboardHistory = append(clipboardHistory, loadedHistory...)
// 	log.Printf("Loaded %d items from history.", len(clipboardHistory))

// 	go connection()
// 	signalChan := make(chan os.Signal, 1)
// 	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

// 	stopWatching := make(chan struct{})

// 	go func() {
// 		cmd := exec.Command("wl-paste", "--watch", "wl-paste")

// 		stdoutPipe, err := cmd.StdoutPipe()
// 		if err != nil {
// 			log.Printf("ERROR: Failed to get stdout pipe for wl-paste: %v", err)
// 			return
// 		}
// 		stderrPipe, err := cmd.StderrPipe()
// 		if err != nil {
// 			log.Printf("ERROR: Failed to get stderr pipe for wl-paste: %v", err)
// 			return
// 		}
// 		if err := cmd.Start(); err != nil {
// 			log.Printf("ERROR: Failed to start wl-paste --watch: %v", err)
// 			return
// 		}
// 		log.Println("Successfully started wl-paste --watch. Monitoring clipboard changes...")

// 		go func() {
// 			scanner := bufio.NewScanner(stdoutPipe)
// 			for scanner.Scan() {
// 				content := scanner.Text()
// 				addText(content)
// 			}
// 			if err := scanner.Err(); err != nil {
// 				log.Printf("ERROR: Error reading from wl-paste stdout: %v", err)
// 			}
// 		}()

// 		// Goroutine to read and log any errors/warnings from wl-paste's stderr
// 		go func() {
// 			scanner := bufio.NewScanner(stderrPipe)
// 			for scanner.Scan() {
// 				log.Printf("WL-PASTE STDERR: %s", scanner.Text())
// 			}
// 		}()

// 		select {
// 		case <-stopWatching:
// 			// Received signal to stop watching, terminate wl-paste process
// 			log.Println("Received stop signal for wl-paste watcher. Attempting to kill wl-paste process...")
// 			if err := cmd.Process.Kill(); err != nil { // Use Kill to ensure termination
// 				log.Printf("WARNING: Failed to kill wl-paste process: %v", err)
// 			}
// 		case err := <-func() chan error {
// 			errChan := make(chan error, 1)
// 			go func() {
// 				errChan <- cmd.Wait() // Wait for wl-paste to exit
// 			}()
// 			return errChan
// 		}():
// 			// wl-paste process exited on its own (could be an error or normal termination)
// 			log.Printf("wl-paste --watch process exited unexpectedly: %v", err)
// 		}
// 	}()

// 	// --- Main goroutine waits for OS signal to shut down the daemon ---
// 	s := <-signalChan // Block until an OS signal is received
// 	log.Printf("Stopping daemon process due to signal: %v", s)

// 	// --- Initiate graceful shutdown sequence ---
// 	close(stopWatching) // Signal the wl-paste watcher goroutine to stop

// 	// Give the wl-paste goroutine a moment to clean up
// 	time.Sleep(500 * time.Millisecond)

// 	// Save clipboard history before exiting
// 	saveClipboardHistory()

// 	log.Println("Daemon process stopped cleanly.")
// }