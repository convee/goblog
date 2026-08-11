package main

import (
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	"github.com/convee/goblog/conf"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/model"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

var (
	configPath = pflag.StringP("config", "c", conf.DefaultPath(), "config file path")
	emailFlag  = pflag.String("email", "", "admin email address")
)

func main() {
	pflag.Parse()
	email := strings.TrimSpace(*emailFlag)
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email || len(email) > 50 {
		log.Fatal("provide a valid --email address")
	}

	password := os.Getenv("GOBLOG_ADMIN_PASSWORD")
	if password == "" {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			log.Fatal("set GOBLOG_ADMIN_PASSWORD or run from an interactive terminal")
		}
		fmt.Print("Admin password: ")
		passwordBytes, readErr := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if readErr != nil {
			log.Fatalf("read password: %v", readErr)
		}
		password = string(passwordBytes)
	}
	if len([]rune(password)) < 12 {
		log.Fatal("admin password must contain at least 12 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}
	cfg := conf.Init(*configPath)
	db := daos.Init(&cfg.Mysql)
	defer db.Close()
	if _, err := daos.AddUser(model.User{Email: email, Password: string(hash)}); err != nil {
		log.Fatalf("create admin: %v", err)
	}
	fmt.Printf("Admin %s created.\n", email)
}
