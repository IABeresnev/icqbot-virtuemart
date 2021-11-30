Получить номер последнего заказа, способ оплаты, способ доставки
SELECT `virtuemart_order_id` FROM `ce7l3_virtuemart_orders` ORDER BY virtuemart_order_id DESC LIMIT 1

Получить способ оплаты
SELECT `payment_name` FROM `ce7l3_virtuemart_paymentmethods_ru_ru` where `virtuemart_paymentmethod_id` = (SELECT `virtuemart_paymentmethod_id` FROM `ce7l3_virtuemart_orders` ORDER BY virtuemart_order_id DESC LIMIT 1)

Получить способ доставки
SELECT `shipment_name` FROM `ce7l3_virtuemart_shipmentmethods_ru_ru` where `virtuemart_shipmentmethod_id` = (SELECT `virtuemart_shipmentmethod_id` FROM `ce7l3_virtuemart_orders` ORDER BY virtuemart_order_id DESC LIMIT 1)

Получить данные покупателя
SELECT `first_name`,`phone_2`,`address_2`,`email` FROM `ce7l3_virtuemart_order_userinfos` where `virtuemart_order_id` = (SELECT `virtuemart_order_id` FROM `ce7l3_virtuemart_orders` ORDER BY virtuemart_order_id DESC LIMIT 1)

Получить позиции заказа
SELECT `virtuemart_order_item_id`,`order_item_sku`,`order_item_name`,`product_final_price` FROM `ce7l3_virtuemart_order_items` where `virtuemart_order_id` = (SELECT `virtuemart_order_id` FROM `ce7l3_virtuemart_orders` ORDER BY virtuemart_order_id DESC LIMIT 1)

Получить id продукта по артикулу.
SELECT `virtuemart_product_id` FROM `ce7l3_virtuemart_products` WHERE product_sku LIKE "УК-00002290"

SELECT `virtuemart_customfield_id`, `virtuemart_custom_id`, `customfield_value` FROM `ce7l3_virtuemart_product_customfields` WHERE `virtuemart_product_id` = (SELECT `virtuemart_product_id` FROM `ce7l3_virtuemart_products` WHERE product_sku LIKE "УК-00002290")

Проверяем наличие доп полей к товару
SELECT `virtuemart_customfield_id`, `virtuemart_custom_id`, `customfield_value` FROM `ce7l3_virtuemart_product_customfields` WHERE `virtuemart_product_id` = (SELECT `virtuemart_product_id` FROM `ce7l3_virtuemart_products` WHERE product_sku LIKE "УК-00008200")

Если результат пустой создаем поля запросом, используя полученный ид продукта по артикулу, ид доп поля нужного магазина и количество товара на складе в этом магазине.
INSERT INTO `ce7l3_virtuemart_product_customfields`(`virtuemart_product_id`, `virtuemart_custom_id`, `customfield_value`) VALUES (ид товара,7,количество на складе) Маркина
INSERT INTO `ce7l3_virtuemart_product_customfields`(`virtuemart_product_id`, `virtuemart_custom_id`, `customfield_value`) VALUES (9642,8,1) Воробьева
INSERT INTO `ce7l3_virtuemart_product_customfields`(`virtuemart_product_id`, `virtuemart_custom_id`, `customfield_value`) VALUES (9642,10,1) Б Хмельницкого  
INSERT INTO `ce7l3_virtuemart_product_customfields`(`virtuemart_product_id`, `virtuemart_custom_id`, `customfield_value`) VALUES (9642,11,1) Кирова
INSERT INTO `ce7l3_virtuemart_product_customfields`(`virtuemart_product_id`, `virtuemart_custom_id`, `customfield_value`) VALUES (9642,12,1) Куликова
INSERT INTO `ce7l3_virtuemart_product_customfields`(`virtuemart_product_id`, `virtuemart_custom_id`, `customfield_value`) VALUES (9642,13,1) Молдавская
INSERT INTO `ce7l3_virtuemart_product_customfields`(`virtuemart_product_id`, `virtuemart_custom_id`, `customfield_value`) VALUES (9642,14,1) Бабаевского
INSERT INTO `ce7l3_virtuemart_product_customfields`(`virtuemart_product_id`, `virtuemart_custom_id`, `customfield_value`) VALUES (9642,15,1) Dexter
INSERT INTO `ce7l3_virtuemart_product_customfields`(`virtuemart_product_id`, `virtuemart_custom_id`, `customfield_value`) VALUES (9642,16,1) камызяк

Если поля уже были, то используя значение `virtuemart_customfield_id` делаем обновление каждого поля отдельно. Значение этого поля будет уникально для каждой строчки.


Для товаров которых нет в базе
Расшифровка
INSERT INTO `ostatki_po_skladam`(`product_name`, `product_art`,`virtuemart_product_id`, `markina`, `vorobieva`, `xmelnickogo`, `kirova`, `kulikova`, `moldavskaya`, `babaika`, `dexter`, `kamiziak`) VALUES ("название товара в 1с","артикул товара в 1с","ID товара в инетмагазе","склад маркина","склад воробьева","склад хмельницкого","склад кирова","склад куликова","молдавская","склад бабайка","склад декстер","склад камызяк")

Пример
INSERT INTO `ostatki_po_skladam`(`product_name`, `product_art`, `markina`, `vorobieva`, `xmelnickogo`, `kirova`, `kulikova`, `moldavskaya`, `babaika`, `dexter`, `kamiziak`) VALUES ("барахло1","артикул барахла",1,2,3,4,5,6,7,8,9)


Для товаров которые уже есть в базе

UPDATE `ostatki_po_skladam` SET `markina`=[value-5],`vorobieva`=[value-6],`xmelnickogo`=[value-7],`kirova`=[value-8],`kulikova`=[value-9],`moldavskaya`=[value-10],`babaika`=[value-11],`dexter`=[value-12],`kamiziak`=[value-13] WHERE `product_art` like "";

Пример обновления полей. Конструкция после `where` должна быть всегда. Конструкция перед может быть любая, так как необязтаельно сразу на всех складах изменилось количество.
UPDATE `ostatki_po_skladam` SET `markina`=11,`vorobieva`=12,`xmelnickogo`=13,`kirova`=14 WHERE `product_art` like "артикул барахла";   


INSERT INTO `ostatki_po_skladam`(`product_name`,`product_art`, `virtuemart_product_id`, `markina`, `vorobieva`, `xmelnickogo`, `kirova`, `kulikova`, `moldavskaya`, `babaika`, `dexter`, `kamiziak`) VALUES ("названиеТовара","кодТовара","IDтовараВинтернетМагазине",,1,3,5,7,9,11,13,15,17) ON DUPLICATE KEY UPDATE `markina`=1, `vorobieva`=3, `xmelnickogo`=5, `kirova`=7, `kulikova`=9, `moldavskaya`=11, `babaika`=13, `dexter`=15, `kamiziak`=17;