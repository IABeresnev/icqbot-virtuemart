package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"

	botgolang "github.com/mail-ru-im/bot-golang"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		User   string `yaml:"user"`
		Pass   string `yaml:"pass"`
		Dbname string `yaml:"dbname"`
	} `yaml:"database"`
	Bot struct {
		Token  string `yaml:"token"`
		Chatid string `yaml:"chatid"`
	} `yaml:"bot"`
}

type Buyer struct {
	first_name string `json:"first_name"`
	phone_2    string `json:"phone_2"`
	address_2  string `json:"address_2"`
	email      string `json:"email"`
}

type Zakaz struct {
	virtuemart_order_item_id string `json:"virtuemart_order_item_id"`
	order_item_sku           string `json:"order_item_sku"`
	order_item_name          string `json:"order_item_name"`
	product_final_price      string `json:"product_final_price"`
}

type PaymentMethod struct {
	payment_name string `json:"payment_name"`
}

type DeliveryMethod struct {
	shipment_name string `json:"shipment_name"`
}

type ldoid struct {
	lastorder string `json:"lastorder"`
}

type OrStatus struct {
	Ordrstat string `json:"Ordrstat"`
}

var CONNSTR string         // "user:pass@tcp(server:port)/dbname" Connection string to mysql db
var TOKEN string           // bot token
var CHATID string          // chatid
var LASTORDERNUMBER string // Last sent order number
var pokupatel string
var zakaztext string
var pmethod string
var dmethod string

func main() {

	var cfg Config
	readFile(&cfg)
	CONNSTR = cfg.Database.User + ":" + cfg.Database.Pass + "@tcp(" + cfg.Server.Host + ":" + cfg.Server.Port + ")/" + cfg.Database.Dbname
	TOKEN = cfg.Bot.Token
	CHATID = cfg.Bot.Chatid

	var ldo, _ = strconv.Atoi(getLastDoneOrderID(CONNSTR))
	var non, _ = strconv.Atoi(getOrderIDForWork(CONNSTR))

	if ldo != non {
		fmt.Println("Let's get id done")
		if (non - ldo) > 1 {
			fmt.Println("More than one order awaits")
			for i := ldo + 1; i <= non; i++ { // +1 для отмены обработки текущего заказа, появлялись дубли.
				if CheckOrderIDStatus(CONNSTR, strconv.Itoa(i)) {
					sendMoreMessage(CONNSTR, TOKEN, CHATID, strconv.Itoa(i))
				} else {
					continue
				}

			}
		} else {
			fmt.Println("One order awaits")
			sendMessage(CONNSTR, TOKEN, CHATID)
		}
	} else {
		fmt.Println("Nothing to be done")
	}

	wrightLastDoneOrderID(CONNSTR, strconv.Itoa(non))
	//sendMessage(CONNSTR, TOKEN, CHATID)
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func sendMessage(fCONNSTR, fTOKEN, fCHATID string) {
	db, err := sql.Open("mysql", fCONNSTR)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// ПОКУПАТЕЛЬ
	resultsPokupatel, err := db.Query("SELECT `first_name`,`phone_2`,`address_2`,`email` FROM `ce7l3_virtuemart_order_userinfos` where `virtuemart_order_id` = (SELECT `virtuemart_order_id` FROM `ce7l3_virtuemart_orders` WHERE `order_status` not like 'P' ORDER BY virtuemart_order_id DESC LIMIT 1)")
	if err != nil {
		panic(err.Error())
	}
	for resultsPokupatel.Next() {
		var buyer Buyer
		err = resultsPokupatel.Scan(&buyer.first_name, &buyer.phone_2, &buyer.address_2, &buyer.email)
		if err != nil {
			panic(err.Error())
		}

		pokupatel = "ФИО -- " + buyer.first_name + " \n" +
			"Телефон -- " + buyer.phone_2 + " \n" +
			"Адрес -- " + buyer.address_2 + " \n" +
			"Электронная почта -- " + buyer.email
	}

	// ПОЗИЦИИ ИЗ ЗАКАЗА
	resultsZakaz, err := db.Query("SELECT `virtuemart_order_item_id`,`order_item_sku`,`order_item_name`,`product_final_price` FROM `ce7l3_virtuemart_order_items` where `virtuemart_order_id` = (SELECT `virtuemart_order_id` FROM `ce7l3_virtuemart_orders` WHERE `order_status` not like 'P' ORDER BY virtuemart_order_id DESC LIMIT 1)")
	if err != nil {
		panic(err.Error())
	}
	var poz int
	poz = 0
	for resultsZakaz.Next() {
		var cart Zakaz
		err = resultsZakaz.Scan(&cart.virtuemart_order_item_id, &cart.order_item_sku, &cart.order_item_name, &cart.product_final_price)
		if err != nil {
			panic(err.Error())
		}
		poz++
		zakaztext = zakaztext + strconv.Itoa(poz) + ") *Название: " + cart.order_item_name + "" + "* *Артикул: " + cart.order_item_sku + "" + "* *Цена на сайте: " + (cart.product_final_price)[:len(cart.product_final_price)-6] + "руб.*" + " \n"
	}

	// СПОСОБ ОПЛАТЫ
	resultsPay, err := db.Query("SELECT `payment_name` FROM `ce7l3_virtuemart_paymentmethods_ru_ru` where `virtuemart_paymentmethod_id` = (SELECT `virtuemart_paymentmethod_id` FROM `ce7l3_virtuemart_orders` WHERE `order_status` not like 'P' ORDER BY virtuemart_order_id DESC LIMIT 1)")
	if err != nil {
		panic(err.Error())
	}
	for resultsPay.Next() {
		var payment PaymentMethod
		err = resultsPay.Scan(&payment.payment_name)
		if err != nil {
			panic(err.Error())
		}
		pmethod = "Способ оплаты: " + payment.payment_name + " \n"
	}

	// СПОСОБ ДОСТАВКИ
	resultsShipment, err := db.Query("SELECT `shipment_name` FROM `ce7l3_virtuemart_shipmentmethods_ru_ru` where `virtuemart_shipmentmethod_id` = (SELECT `virtuemart_shipmentmethod_id` FROM `ce7l3_virtuemart_orders` WHERE `order_status` not like 'P' ORDER BY virtuemart_order_id DESC LIMIT 1)")
	if err != nil {
		panic(err.Error())
	}
	for resultsShipment.Next() {
		var shipment DeliveryMethod
		err = resultsShipment.Scan(&shipment.shipment_name)
		if err != nil {
			panic(err.Error())
		}
		dmethod = "Способ доставки: " + shipment.shipment_name + " \n"
	}

	// BOOOOTTTT
	bot, err := botgolang.NewBot(fTOKEN, botgolang.BotDebug(true))
	if err != nil {
		log.Fatalf("cannot connect to bot: %s", err)
	}

	log.Println(bot.Info)

	message := bot.NewTextMessage(fCHATID, "*********************НОВЫЙ ЗАКАЗ НА САЙТЕ №"+getOrderIDForWork(CONNSTR)+"********************* \n"+pokupatel+"\n\n"+pmethod+" "+dmethod+"\n"+zakaztext)
	if err = message.Send(); err != nil {
		log.Fatalf("failed to send message: %s", err)
	}
}

func sendMoreMessage(fCONNSTR, fTOKEN, fCHATID, fORDERID string) {
	db, err := sql.Open("mysql", fCONNSTR)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// ПОКУПАТЕЛЬ
	resultsPokupatel, err := db.Query("SELECT `first_name`,`phone_2`,`address_2`,`email` FROM `ce7l3_virtuemart_order_userinfos` where `virtuemart_order_id` = " + fORDERID)
	if err != nil {
		panic(err.Error())
	}
	for resultsPokupatel.Next() {
		var buyer Buyer
		err = resultsPokupatel.Scan(&buyer.first_name, &buyer.phone_2, &buyer.address_2, &buyer.email)
		if err != nil {
			panic(err.Error())
		}

		pokupatel = "ФИО -- " + buyer.first_name + " \n" +
			"Телефон -- " + buyer.phone_2 + " \n" +
			"Адрес -- " + buyer.address_2 + " \n" +
			"Электронная почта -- " + buyer.email
	}

	// ПОЗИЦИИ ИЗ ЗАКАЗА
	resultsZakaz, err := db.Query("SELECT `virtuemart_order_item_id`,`order_item_sku`,`order_item_name`,`product_final_price` FROM `ce7l3_virtuemart_order_items` where `virtuemart_order_id` = " + fORDERID)
	if err != nil {
		panic(err.Error())
	}
	var poz int
	poz = 0
	zakaztext = ""
	for resultsZakaz.Next() {
		var cart Zakaz
		err = resultsZakaz.Scan(&cart.virtuemart_order_item_id, &cart.order_item_sku, &cart.order_item_name, &cart.product_final_price)
		if err != nil {
			panic(err.Error())
		}
		poz++
		zakaztext = zakaztext + strconv.Itoa(poz) + ") *Название: " + cart.order_item_name + "" + "* *Артикул: " + cart.order_item_sku + "" + "* *Цена на сайте: " + (cart.product_final_price)[:len(cart.product_final_price)-6] + "руб.*" + " \n"
	}

	// СПОСОБ ОПЛАТЫ
	resultsPay, err := db.Query("SELECT `payment_name` FROM `ce7l3_virtuemart_paymentmethods_ru_ru` where `virtuemart_paymentmethod_id` = (SELECT `virtuemart_paymentmethod_id` FROM `ce7l3_virtuemart_orders` WHERE virtuemart_order_id = " + fORDERID + ")")
	if err != nil {
		panic(err.Error())
	}
	for resultsPay.Next() {
		var payment PaymentMethod
		err = resultsPay.Scan(&payment.payment_name)
		if err != nil {
			panic(err.Error())
		}
		pmethod = "Способ оплаты: " + payment.payment_name + " \n"
	}

	// СПОСОБ ДОСТАВКИ
	resultsShipment, err := db.Query("SELECT `shipment_name` FROM `ce7l3_virtuemart_shipmentmethods_ru_ru` where `virtuemart_shipmentmethod_id` = (SELECT `virtuemart_shipmentmethod_id` FROM `ce7l3_virtuemart_orders` WHERE virtuemart_order_id = " + fORDERID + ")")
	if err != nil {
		panic(err.Error())
	}
	for resultsShipment.Next() {
		var shipment DeliveryMethod
		err = resultsShipment.Scan(&shipment.shipment_name)
		if err != nil {
			panic(err.Error())
		}
		dmethod = "Способ доставки: " + shipment.shipment_name + " \n"
	}

	// BOOOOTTTT
	bot, err := botgolang.NewBot(fTOKEN, botgolang.BotDebug(true))
	if err != nil {
		log.Fatalf("cannot connect to bot: %s", err)
	}

	log.Println(bot.Info)

	message := bot.NewTextMessage(fCHATID, "*********************НОВЫЙ ЗАКАЗ НА САЙТЕ №"+fORDERID+"********************* \n"+pokupatel+"\n\n"+pmethod+" "+dmethod+"\n"+zakaztext)
	if err = message.Send(); err != nil {
		log.Fatalf("failed to send message: %s", err)
	}
}

func getLastDoneOrderID(fCONNSTR string) string {
	db, err := sql.Open("mysql", fCONNSTR)
	var lordernumber string
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// getLastDoneOrderID
	resultsLastDoneOrder, err := db.Query("SELECT `lastorder` FROM `virt_order_check` where `id` = 1")
	if err != nil {
		panic(err.Error())
	}
	for resultsLastDoneOrder.Next() {
		var lnumberdoneorder ldoid
		err = resultsLastDoneOrder.Scan(&lnumberdoneorder.lastorder)
		if err != nil {
			panic(err.Error())
		}

		lordernumber = lnumberdoneorder.lastorder
	}
	return lordernumber
}

func getOrderIDForWork(fCONNSTR string) string {
	db, err := sql.Open("mysql", fCONNSTR)
	var IDForWork string
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// getLastDoneOrderID
	resultsNewOrder, err := db.Query("SELECT `virtuemart_order_id` FROM `ce7l3_virtuemart_orders` WHERE `order_status` not like 'P' ORDER BY virtuemart_order_id DESC LIMIT 1")
	if err != nil {
		panic(err.Error())
	}
	for resultsNewOrder.Next() {
		var lnumberdoneorder ldoid
		err = resultsNewOrder.Scan(&lnumberdoneorder.lastorder)
		if err != nil {
			panic(err.Error())
		}

		IDForWork = lnumberdoneorder.lastorder
	}
	return IDForWork
}

func CheckOrderIDStatus(fCONNSTR string, OrderID string) bool {
	db, err := sql.Open("mysql", fCONNSTR)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// CheckOrderIDStatus
	resultsNewOrder, err := db.Query("SELECT `order_status` FROM `ce7l3_virtuemart_orders` WHERE `virtuemart_order_id` =" + OrderID)
	if err != nil {
		panic(err.Error())
	}
	for resultsNewOrder.Next() {
		var thisOrderStatus OrStatus
		err = resultsNewOrder.Scan(&thisOrderStatus.Ordrstat)
		if err != nil {
			panic(err.Error())
		}

		if thisOrderStatus.Ordrstat != "P" {
			return true
		} else {
			return false
		}
	}

	return false
}

func wrightLastDoneOrderID(fCONNSTR, lastDoneOrderNumber string) {
	db, err := sql.Open("mysql", fCONNSTR)
	if err != nil {
		panic(err.Error())
	}

	insForm, err := db.Prepare("UPDATE virt_order_check SET lastorder=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(lastDoneOrderNumber, 1)
	defer db.Close()
}
