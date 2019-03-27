package main

import (
	"github.com/labstack/gommon/log"
	"net"
	"os"
	"plugin"
)

var (
	config Config
	logger *log.Logger
)

func main() {
	var (
		store Connector
	)
	logger = log.New("-")

	if len(os.Args) == 2 {
		if err := config.Load(os.Args[1]); err != nil {
			logger.Fatalf("Ошибка парсинга конфига: %v", err)
		}
	} else {
		logger.Fatalf("Не задан путь до конфига")
	}
	logger.SetLevel(config.GetLogLevel())

	if config.Store != nil {
		plug, err := plugin.Open(config.Store["plugin"])
		if err != nil {
			logger.Fatalf("Не удалость загрузить плагин хранилища: %v", err)
		}

		connector, err := plug.Lookup("Connector")
		if err != nil {
			logger.Fatalf("Не удалось загрузить коннектор: %v", err)
		}

		store = connector.(Connector)
	} else {
		store = defaultConnector{}
	}

	if err := store.Init(config.Store); err != nil {
		logger.Fatal(err)
	}
	defer store.Close()

	l, err := net.Listen("tcp", config.GetListenAddress())
	if err != nil {
		logger.Fatalf("Не удалось открыть соединение: %v", err)
	}
	defer l.Close()

	logger.Infof("Запущен сервер %s...", config.GetListenAddress())
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Errorf("Ошибка соединения: %v", err)
		} else {
			go handleRecvPkg(conn, store)

		}
	}
}
