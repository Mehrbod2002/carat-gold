package utils

import (
	"context"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var App *firebase.App

func InitializeApp() error {
	ctx := context.Background()
	mainPath, err := os.Getwd()
	if err != nil {
		return err
	}
	opt := option.WithCredentialsFile(filepath.Join(mainPath,
		"carat-b4654-firebase-adminsdk-ozsc4-84e240d33e.json"))
	var errApp error
	App, errApp = firebase.NewApp(ctx, nil, opt)
	if errApp != nil {
		return err
	}
	return nil
}

func GetApp() *firebase.App {
	return App
}
