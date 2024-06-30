package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

func main() {
	// Устанавливаем соединение с Netlink для диагностики сетевых интерфейсов
	conn, err := netlink.Dial(unix.NETLINK_INET_DIAG, nil)
	if err != nil {
		log.Fatalf("Ошибка при установке соединения с Netlink: %v", err)
	}
	defer conn.Close()

	log.Println("Соединение с Netlink успешно установлено.")

	// Формируем запрос на получение информации о TCP сокетах
	req := netlink.Message{
		Header: netlink.Header{
			Type:     unix.SOCK_DIAG_BY_FAMILY,
			Flags:    netlink.Request,
			Sequence: 1,
			PID:      uint32(os.Getpid()),
		},
	}

	// Отправляем запрос и получаем ответ
	msgs, err := conn.Execute(req)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса Netlink: %v", err)
	}

	fmt.Println("msgs", msgs)

	// Обрабатываем ответ и выводим информацию о TCP сокетах

}
