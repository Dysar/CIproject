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
		if gotest(){
			fmt.Println("Test passes")
			build() //build binaries
		} else {
			fmt.Println("Test failed")
		}
		//send binaries to Slack via slackbot mazafaka
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
	//cmdName3 := "git"
	//cmdArgs3 := "init" //git init
	//if _, err := exec.Command(cmdName3, cmdArgs3).Output(); err != nil {
	//	fmt.Fprintln(os.Stderr, "There was an error running git init command in local repository: ", err)
	//}

	cmdName := "git"
	cmdArgs := []string {"clone","--bare","https://github.com/Dysar/CIproject.git", path} //git clone --bare
	if _, err := exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Repository already exits", err)
	}

	//cmdName4 := "git"
	//cmdArgs4 := []string {"pull", "https://github.com/Dysar/CIproject"} //git pull
	//if _, err := exec.Command(cmdName4, cmdArgs4...).Output(); err != nil {
	//	fmt.Fprintln(os.Stderr, "There was an error running git pull command: ", err)
	//	os.Exit(1)
	//}
	//
	//cmdName2 := "git"
	//cmdArgs2 := []string {"fetch", "https://github.com/Dysar/CIproject"} //git fetch
	//if _, err := exec.Command(cmdName2, cmdArgs2...).Output(); err != nil {
	//	fmt.Fprintln(os.Stderr, "There was an error running git fetch command: ", err)
	//	os.Exit(1)
	//}

}
func gitcheckout(hash string){
	cmdName := "git"
	cmdArgs := []string {"reset","--hard", hash} //git checkout
	if err := exec.Command(cmdName, cmdArgs...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git reset command in commit repo: ", err) //error is that you have to cd into that dir, localrepo, and even there is says Could not parse or is not a tree
		os.Exit(1)
	}
}