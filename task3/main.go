package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:id"`
	Email string `json:email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) (err error) {
	id := args["id"]
	operation := args["operation"]
	fileName := args["fileName"]
	item := args["item"]

	if len(operation) == 0 {
		return errors.New("-operation flag has to be specified")
	}
	if len(fileName) == 0 {
		return errors.New("-fileName flag has to be specified")
	}

	switch operation {
	case "list":
		data, _ := ioutil.ReadFile(fileName)
		writer.Write(data)
	case "add":
		if len(item) == 0 {
			return errors.New("-item flag has to be specified")
		}
		err = add(item, fileName)
		return
	case "findByID":
		if len(id) == 0 {
			return errors.New("-id flag has to be specified")
		}
		var user User
		user, err = findById(id, fileName)
		var data []byte
		if user.Id != "" {
			data, _ = json.Marshal(user)
		} else {
			data = []byte("")
		}
		writer.Write(data)
	case "remove":
		if len(id) == 0 {
			return errors.New("-id flag has to be specified")
		}
		err = remove(id, fileName)

	default:
		return errors.New("Operation " + operation + " not allowed!")
	}
	return nil
}

func add(item, fileName string) (err error) {
	var user User
	err = json.Unmarshal([]byte(item), &user)
	if err != nil {
		return
	}
	users, err := getUsersFromFile(fileName)
	if err != nil {
		return
	}

	if userExist(user, users) {
		return errors.New("Item with id " + user.Id + " already exists")
	}
	users = append(users, user)
	err = writeToFile(users, fileName)
	return
}

func writeToFile(users []User, fileName string) (err error) {
	usersForJson, err := json.Marshal(users)
	if err != nil {
		return
	}
	return ioutil.WriteFile(fileName, usersForJson, 0677)
}

func getUsersFromFile(fileName string) (users []User, err error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &users)
	file.Close()
	return
}

func userExist(newUser User, users []User) bool {
	for _, userFromList := range users {
		if userFromList.Id == newUser.Id {
			return true
		}
	}
	return false
}

func findById(id, fileName string) (user User, err error) {
	users, err := getUsersFromFile(fileName)
	if err != nil {
		return
	}
	for _, userFromFile := range users {
		if userFromFile.Id == id {
			user = userFromFile
			return
		}
	}
	return
}

func remove(id, fileName string) (err error) {
	users, err := getUsersFromFile(fileName)
	user, err := findById(id, fileName)

	if err != nil {
		return
	}

	if user.Id != "" {
		var newListUsers []User
		for _, user := range users {
			if user.Id != id {
				newListUsers = append(newListUsers, user)
			}
		}
		return writeToFile(newListUsers, fileName)
	}
	return errors.New("Item with id " + id + " not found")
 }

func parseArgs() Arguments {
	id := flag.String("id", "", "id")
	item := flag.String("item", "", "item")
	operation := flag.String("operation", "", "operation")
	fileName := flag.String("fileName", "", "fileName")
	flag.Parse()
	return Arguments{"id": *id, "item": *item, "operation": *operation, "fileName": *fileName}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
