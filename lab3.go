package main

import (
	"fmt"
	"time"
)

// Token - структура для передачи сообщения
type Token struct {
	Data      string
	Recipient int
	TTL       int
}

func main() {
	var N int
	fmt.Print("Введите количество узлов: ")
	fmt.Scan(&N)

	// Создание каналов для каждого узла
	channels := make([]chan Token, N)
	for i := range channels {
		channels[i] = make(chan Token)
	}

	// Инициализация узлов
	for i := 0; i < N; i++ {
		next := (i + 1) % N // следующий узел в кольце
		go node(i, channels[i], channels[next])
	}

	// Отправка токена первому узлу
	go func() {
		channels[0] <- Token{Data: "Привет, Token Ring!", Recipient: N / 2, TTL: N + 1}
	}()

	// Даем время на выполнение передачи
	time.Sleep(5 * time.Second)
}

// Функция узла
func node(id int, in <-chan Token, out chan<- Token) {
	for {
		token := <-in // Получение токена от предыдущего узла
		if token.TTL <= 0 {
			continue // Если TTL истек, пропускаем токен
		}

		// Если этот узел является получателем
		if id == token.Recipient {
			fmt.Printf("Узел %d получил сообщение: %s\n", id, token.Data)
			continue // После получения сообщения узел продолжает работу
		}

		// Уменьшаем TTL и передаем токен следующему узлу
		token.TTL--
		out <- token
	}
}
