package main

import (
	"fmt"
	"github.com/phayes/hookserve/hookserve"
	"github.com/nlopes/slack"
	"os"
	"os/exec"
	"time"
	"log"
)

func main() {
	server := hookserve.NewServer()
	server.Port = 8888
	server.Path = "/postreceive"
	server.Secret = "supersecretcode"
	server.GoListenAndServe()

	api := slack.New("xoxb-107838516693-fPRh1G5nMAG69F7PD4Lw7evq")
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	gitpreparation()

	// Everytime the server receives a webhook event, print the results
	for event := range server.Events {
		fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
		path := "/home/ubuntu"
		if err := os.Chdir(path); err != nil {
			panic(err)
		}
		cmdName := "git"
		cmdArgs := []string {"clone","/home/ubuntu/localrepo.git", "localrepo"} //git clone --bare, git fetch, git clone, git checkout
		if _, err := exec.Command(cmdName, cmdArgs...).Output(); err != nil {
			fmt.Fprintln(os.Stderr, "There was an error running git clone command: ", err)
			os.Exit(1)
		} //download latest files
		err := os.Rename("localrepo",event.Commit)
		if err != nil {
			panic(err)
		}
		gitcheckout(event.Commit)
		gotest() //run the test in that repo go test in that repo
		rtm := api.NewRTM()
		go rtm.ManageConnection()
		if gotest(){
			msg := fmt.Sprintf("Tests passed for %d commit", event.Commit)
			slackbot(rtm, msg)
			os.Exit(0)
		} else {
			msg := fmt.Sprintf("Tests failed for %d commit", event.Commit)
			slackbot(rtm, msg)
			os.Exit(1)
		}
		//send binaries to Slack via slackbot mazafaka
	}
}
func gotest()(bool){
	cmdName := "go"
	cmdArgs := "test"
	if _, err := exec.Command(cmdName, cmdArgs).CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running tests: ", err)
		return false
	}
	return true
}
func build(){
	cmdName := "go"
	cmdArgs := []string {"build","main.go"}
	if _, err := exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Build failed: ", err)
		os.Exit(1)
	}
	fmt.Println("Binaries were build")
}
func gitpreparation(){
	path := "/home/ubuntu/localrepo.git"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 755)
	}

	if err := os.Chdir(path); err != nil {
		panic(err)
	}

	cmdName := "git"
	cmdArgs := []string {"clone","--bare","https://github.com/Dysar/CIproject.git", path} //git clone --bare
	if _, err := exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Repository already exits", err)
	}

}
func gitcheckout(hash string){
	fmt.Println(time.Now)
	if err := os.Chdir(hash); err != nil {
		panic(err)
	}
	cmdName4 := "git"
	cmdArgs4 := []string {"pull", "https://github.com/Dysar/CIproject"} //to refresh git log, WITHOUT THIS CHECKOUT FAILS
	if _, err := exec.Command(cmdName4, cmdArgs4...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git pull command: ", err)
		os.Exit(1)
	}
	cmdName := "git"
	cmdArgs := []string {"checkout", hash} //git checkout
	if err := exec.Command(cmdName, cmdArgs...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git checkout command in commit repo: ", err) //error is that you have to cd into that dir, localrepo, and even there is says Could not parse or is not a tree
		os.Exit(1)
	}
}

func slackbot(rtm *slack.RTM, message string) {
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			//case *slack.HelloEvent:
			//	rtm.SendMessage(rtm.NewOutgoingMessage("Hello event received", "D351BB3EC"))

			case *slack.ConnectedEvent:
				fmt.Println("Infos:", ev.Info)
				fmt.Println("Connection counter:", ev.ConnectionCount)
				rtm.SendMessage(rtm.NewOutgoingMessage(message, "D351BB3EC"))
				return
			// Replace #general with your Channel ID
			}
		}
	}
	return
}