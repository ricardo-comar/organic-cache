package main

import (
	"fmt"
	"time"
)

func main() {
	// Defina um timeout de 3 segundos
	timeout := 3 * time.Second

	// Crie um canal para sinalizar quando a tarefa estiver concluída
	taskCompleted := make(chan bool)

	// Inicie a função da tarefa em uma goroutine
	go performTask(taskCompleted)

	// Use select para esperar pela conclusão da tarefa ou pelo timeout
	select {
	case <-taskCompleted:
		fmt.Println("Tarefa concluída com sucesso.")
	case <-time.After(timeout):
		fmt.Println("Timeout! A tarefa demorou muito para ser concluída.")
	}
}

func performTask(done chan bool) {
	// Simule uma tarefa que leva 2 segundos para ser concluída
	time.Sleep(2 * time.Second)

	// Quando a tarefa estiver concluída, sinalize através do canal
	done <- true
}
