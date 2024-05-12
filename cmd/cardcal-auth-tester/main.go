package main

import (
	"cardcaldav"
	"context"
	"flag"
	"fmt"
	"github.com/charmbracelet/log"
	"golang.org/x/term"
	"os"
)

var Logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportTimestamp: true,
	ReportCaller:    true,
})

func main() {
	var un, dbStr string
	var debugMode bool
	flag.StringVar(&un, "un", "", "username of user to authenticate")
	flag.StringVar(&dbStr, "db", "", "Connection string for the database - user:password@tcp(127.0.0.1:3306)/dbname")
	flag.BoolVar(&debugMode, "debug", false, "enable debug logging")
	flag.Parse()

	if debugMode {
		Logger.SetLevel(log.DebugLevel)
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		Logger.Fatal("No terminal input found")
	}

	auth := cardcaldav.NewAuth(dbStr, Logger)

	_, _ = fmt.Fprint(os.Stderr, "Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		Logger.Fatal("ReadPassword()", "err", err)
	}
	_, _ = fmt.Fprintln(os.Stderr)
	Logger.Info("Verifying password...")

	err = auth.ValidateCredentials(context.Background(), un, string(password))
	if err != nil {
		Logger.Fatal("auth.ValidateCredentials()", "err", err)
	}

	Logger.Info("Password verified")
}
