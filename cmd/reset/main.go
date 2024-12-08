package main

import (
	"fmt"
	"os"

	"github.com/Zambozoo/baby-names/src/db"
	"github.com/Zambozoo/baby-names/src/terminal"
)

const (
	userDBEnv           = "USER_DB"
	userDBPathPrompt    = "What is the path to the database file? (can be new file):\n"
	resetIdentityPrompt = "What is your name? (deletes both you and your partner's account):\n"
)

func main() {
	fmt.Printf("Welcome to Baby Names! At any time, press Ctrl+C to exit.\n")

	tp := terminal.NewTerminalPrompter(os.Stdin, os.Stdout)
	userDB := getUserDB(tp)

	username := tp.Prompt(stringEmpty, resetIdentityPrompt)
	user, err := userDB.DeleteUsers(username)
	check(err)
	fmt.Printf("Deleted users [%s,%s]\n", user.Username, user.PartnerUsername)
}

func getUserDB(tp *terminal.TerminalPrompter) *db.UserDB {
	dbFilePath := os.Getenv(userDBEnv)
	if dbFilePath == "" {
		printCurrentDirectory()
		dbFilePath = tp.Prompt(stringEmpty, userDBPathPrompt)
	}

	userDB, err := db.NewUserDB(dbFilePath)
	check(err)

	return userDB
}

func printCurrentDirectory() {
	dir, err := os.Getwd()
	check(err)
	fmt.Printf("Current directory: %s\n", dir)
}

func stringEmpty(s string) bool {
	return len(s) != 0
}

func yesOrNo(s string) bool {
	return s == "y\n" || s == "n\n"
}

func check(err error) {
	if err != nil {
		exitErr(err.Error())
	}
}

func exitErr(msg string) {
	fmt.Println(msg)
	os.Exit(-1)
}
