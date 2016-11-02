package main

import (
	"fmt"
	"github.com/phayes/hookserve/hookserve"
	"os"
	"os/exec"
)

func main() {
	server := hookserve.NewServer()
	server.Port = 8888
	server.Path = "/postreceive"
	server.Secret = "supersecretcode"
	server.GoListenAndServe()

	gitpreparation()

	// Everytime the server receives a webhook event, print the results
	for event := range server.Events {
		fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
		if err := os.Chdir("/home/ubuntu"); err != nil {
			panic(err)
		}
		err := os.Mkdir(event.Commit, 755) //create clone repo for this commit
		if err != nil {
			panic(err)
		}
		if err = os.Chdir(event.Commit); err != nil {
			panic(err)
		}
		cmdName := "git"
		cmdArgs := []string {"clone","/home/ubuntu/localrepo.git"} //git clone --bare, git fetch, git clone, git checkout
		var out []byte
		if out, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
			fmt.Fprintln(os.Stderr, "There was an error running git clone command: ", err)
			fmt.Fprintln(os.Stderr,"This bug", out)
			os.Exit(1)
		} //download latest files
		gitcheckout(event.Commit)
		gotest() //run the test in that repo go test in that repo
		if err != nil {
			panic(err) //check the test status
		}
		build() //build binaries
		//send binaries to Slack via slackbot
	}
}
func gotest()(bool){
	cmdName := "go"
	cmdArgs := "test"
	if _, err := exec.Command(cmdName, cmdArgs).Output(); err != nil {
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
}
func gitpreparation(){
	path := "/home/ubuntu/localrepo.git"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 755)
	}

	if err := os.Chdir(path); err != nil {
		panic(err)
	}
	cmdName3 := "git"
	cmdArgs3 := "init" //git init
	if _, err := exec.Command(cmdName3, cmdArgs3).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git init command in local repository: ", err)
	}

	cmdName := "git"
	cmdArgs := []string {"clone","--bare","https://github.com/Dysar/CIproject.git"} //git clone --bare
	if _, err := exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Repository already exits", err)
	}

	cmdName2 := "git"
	cmdArgs2 := []string {"fetch", "https://github.com/Dysar/CIproject"} //git fetch
	if _, err := exec.Command(cmdName2, cmdArgs2...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git fetch command: ", err)
		os.Exit(1)
	}

}
func gitcheckout(hash string){
	cmdName := "git"
	cmdArgs := []string {"reset","--hard", hash} //git checkout
	if err := exec.Command(cmdName, cmdArgs...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git reset command in commit repo: ", err) //error is that you have to cd into that dir, localrepo, and even there is says Could not parse or is not a tree
		os.Exit(1)
	}
}