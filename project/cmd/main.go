package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/configor"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"os"
	"project/api"
	"project/cmd/my_service"
	"project/cmd/my_service/cache"
	"project/cmd/my_service/repository"
	"project/configs"
	"strconv"
	"sync"
	"time"
)

func main() {
	cfg := getConfig("configs/config.yaml")

	configureLogger(cfg)

	db, err := repository.ConnectPostgresDB(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	repo := repository.ItemPostgresSQLItemRepository(db)

	if err != nil {
		log.Fatal("can't get count of from repo %w", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	cacheItems, err := repo.GetItems(ctx)
	if err != nil {
		log.Fatal("service can't get items from repo :%w ", err)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	serviceCache := cache.NewMapCache(cacheItems)

	itemService := my_service.New(repo, serviceCache)

	httpHandler := configureHttpRouter(itemService, cfg)
	go func() {
		runService(cfg, httpHandler)
	}()

	natsMessagesCh := configureNats(cfg)
	go func() {
		runProcessingNatsMessages(natsMessagesCh, itemService)
	}()

	wg.Wait()
}

func configureLogger(cfg configs.Config) {
	f, err := os.OpenFile(cfg.GetLogPath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("can't open log file %s", cfg.GetLogPath())
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Service Logger: ")
}

func runProcessingNatsMessages(msgCh chan *nats.Msg, itemService *my_service.ItemService) {
	for {
		select {
		case newMessage := <-msgCh:
			{
				go worker(newMessage.Data, itemService)
			}
		}
	}
}

func worker(jsonData []byte, itemService *my_service.ItemService) {
	var item api.Item
	err := json.Unmarshal(jsonData, &item)
	if err != nil {
		log.Printf("can't unmarshal message from nats channel %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err = itemService.AddItem(ctx, item)

	if err != nil {
		fmt.Printf("can't add item to service %s", err)
	}
}

func runService(cfg configs.Config, handler http.Handler) {
	addr := ":" + strconv.Itoa(cfg.GetServerPort())
	log.Fatal(http.ListenAndServe(addr, handler))
}

func configureNats(cfg configs.Config) chan *nats.Msg {
	nc, err := nats.Connect(cfg.GetNATSUrl())
	if err != nil {
		log.Fatal("can't connect to nats server")
	}
	ch := make(chan *nats.Msg, 64)
	const subject = "json_payload"
	_, err = nc.ChanSubscribe(subject, ch)
	if err != nil {
		log.Fatalf("nats can't subscribe to subject %s", subject)
	}
	return ch
}

func configureHttpRouter(itemService *my_service.ItemService, cfg configs.Config) http.Handler {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := my_service.NewServer(itemService)

	router.HandleFunc("/task", server.CreateItemHandler).Methods("POST")
	router.HandleFunc("/task/{uuid}", server.GetItemHandler).Methods("GET")

	handler := handlers.RecoveryHandler()(router)
	credentials := handlers.AllowCredentials()

	headersOk := handlers.AllowedHeaders([]string{
		"Access-Control-Allow-Headers", "Origin", "X-Requested-With", "Content-Type", "Accept"})
	frontendPort := strconv.Itoa(cfg.GetFrontendPort())
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:" + frontendPort})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	return handlers.CORS(headersOk, originsOk, methodsOk, credentials)(handler)
}

func getConfig(path string) configs.Config {
	var cfg configs.Config
	err := configor.Load(&cfg, path)
	if err != nil {
		log.Panic("can't open config.yaml %w", err)
	}
	return cfg
}
