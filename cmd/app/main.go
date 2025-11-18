package main

import (
	"context"
	"log"

	"github.com/UdinSemen/golang-test-task/internal/app"
)

// @title			Сервис для тестового задания
// @version		1.0.0
// @description	Данный сервис принимает запросы на добавление числа и возвращает отсортированный список чисел
// @BasePath		/v1
func main() {
	ctx := context.Background()

	myApp, err := app.New(ctx)
	if err != nil {
		log.Fatal("create app", err)
	}

	if err = myApp.Run(); err != nil {
		log.Fatal("run app", err)
	}
}
