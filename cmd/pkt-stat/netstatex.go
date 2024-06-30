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
			Flags:    netlink.Request, // Базовый запрос без дополнительных флагов
			Sequence: 1,
			PID:      uint32(os.Getpid()),
		},
	}

	// Формируем пустой атрибут запроса
	emptyAttr := netlink.Attribute{
		Type: myconst.INET_DIAG_REQ_BYTECODE, // Здесь укажите нужный тип атрибута
		Data: []byte{},                       // Пустые данные атрибута
	}

	// Создаем атрибуты запроса
	attrs, err := netlink.MarshalAttributes([]netlink.Attribute{emptyAttr})
	if err != nil {
		log.Fatalf("Ошибка при формировании атрибутов запроса: %v", err)
	}

	// Устанавливаем атрибуты запроса в сообщение Netlink
	req.Data = attrs

	// Отправляем запрос и получаем ответ
	msgs, err := conn.Execute(req)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса Netlink: %v", err)
	}

	fmt.Println("msgs", msgs)

	// Обрабатываем ответ и выводим информацию о TCP сокетах
}
