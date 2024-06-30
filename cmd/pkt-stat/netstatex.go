package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Gazmasater/myconst"
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
			Flags:    netlink.Request | netlink.Dump,
			Sequence: 1,
			PID:      uint32(os.Getpid()),
		},
	}

	// Создаем атрибуты запроса
	//myconst.INET_DIAG_REQ_BYTECODE   взял с пакета С
	attrs, err := netlink.MarshalAttributes([]netlink.Attribute{
		{Type: myconst.INET_DIAG_REQ_BYTECODE, Data: []byte{}},
	})

	if err != nil {
		log.Fatalf("Ошибка при формировании атрибутов запроса: %v", err)
	}

	// Устанавливаем атрибуты запроса в сообщение Netlink
	req.Data = attrs

	// Выводим содержимое req.Data
	fmt.Printf("Содержимое атрибутов запроса:\n")
	fmt.Printf("%v\n", req.Data)

	// Отправляем запрос и получаем ответ
	msgs, err := conn.Execute(req)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса Netlink: %v", err)
	}

	log.Println("Получена информация о TCP соединениях:")

	// Обрабатываем ответ и выводим информацию о TCP сокетах
	for _, msg := range msgs {
		fmt.Printf("Сообщение Netlink:\n")
		//fmt.Printf("Длина данных: %d\n", msg.Header.Len)
		fmt.Printf("Тип сообщения: %d\n", msg.Header.Type)
		fmt.Printf("Последовательность: %d\n", msg.Header.Sequence)
		fmt.Printf("Идентификатор процесса: %d\n", msg.Header.PID)
		fmt.Printf("Данные: %v\n", msg.Data)
	}
}
