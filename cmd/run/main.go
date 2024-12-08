package main

import (
	"fmt"
	"os"

	"github.com/Zambozoo/baby-names/src/db"
	"github.com/Zambozoo/baby-names/src/terminal"
)

const (
	userDBEnv   = "USER_DB"
	nameFileEnv = "NAMES_FILE"

	userDBPathPrompt    = "What is the path to the database file? (can be new file):\n"
	namesFilePathPrompt = "What is the path to the names file? (newline delimited):\n"

	identityPrompt = "What is your name? (creates your account):\n"
	partnerPrompt  = "Who is your partner? (creates their account):\n"

	likeNamePrompt   = "Do you like the name %q? (y/n):\n"
	nameMatchMessage = "You and your partner matched on the name %q!\n"
)

func main() {
	fmt.Printf("Welcome to Baby Names! At any time, press Ctrl+C to exit.\n")

	tp := terminal.NewTerminalPrompter(os.Stdin, os.Stdout)
	userDB := getUserDB(tp)
	names := getNames(tp)
	user, partnerUser := getUsers(tp, userDB)

	// Print matches
	if matched := user.Matched(partnerUser); len(matched) > 0 {
		fmt.Printf("You and your partner have matched the following names: %v\n", matched)
	} else {
		fmt.Println("You and your partner have not matched any names.")
	}

	// Iterate through names partner has liked
	for name := range partnerUser.LikedNames {
		if _, ok := user.LikedNames[name]; ok {
			continue
		} else if _, ok := user.DislikedNames[name]; ok {
			continue
		}

		if tp.Prompt(yesOrNo, fmt.Sprintf(likeNamePrompt, name)) == "y" {
			user.LikedNames[name] = struct{}{}
			fmt.Printf(nameMatchMessage, name)
		} else {
			user.DislikedNames[name] = struct{}{}
		}

		check(userDB.UpdateUser(user))
	}

	// Iterate through new names
	for _, name := range names {
		if _, ok := user.LikedNames[name]; ok {
			continue
		} else if _, ok := user.DislikedNames[name]; ok {
			continue
		}

		if tp.Prompt(yesOrNo, fmt.Sprintf(likeNamePrompt, name)) == "y" {
			user.LikedNames[name] = struct{}{}
		} else {
			user.DislikedNames[name] = struct{}{}
		}

		check(userDB.UpdateUser(user))
	}
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

func getNames(tp *terminal.TerminalPrompter) []string {
	namesFilePath := os.Getenv(nameFileEnv)
	if namesFilePath == "" {
		printCurrentDirectory()
		namesFilePath = tp.Prompt(stringEmpty, namesFilePathPrompt)
	}

	names, err := db.ReadNames(namesFilePath)
	check(err)

	return names
}

func printCurrentDirectory() {
	dir, err := os.Getwd()
	check(err)
	fmt.Printf("Current directory: %s\n", dir)
}

func getUsers(tp *terminal.TerminalPrompter, userDB *db.UserDB) (*db.User, *db.User) {
	username := tp.Prompt(stringEmpty, identityPrompt)
	user, err := userDB.GetUser(username)
	check(err)

	var partnerUser *db.User
	if user == nil {
		partnerUsername := tp.Prompt(stringEmpty, partnerPrompt)
		partner, err := userDB.GetUser(partnerUsername)
		check(err)
		if partner != nil {
			exitErr("Cannot link new user to existing partner.\n")
		}

		user, partnerUser, err = userDB.CreateUsers(username, partnerUsername)
		check(err)
	} else {
		partnerUser, err = userDB.GetUser(user.PartnerUsername)
		check(err)
		if partnerUser == nil {
			exitErr("Missing partner from UserDB.\n")
		}
	}

	return user, partnerUser
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
