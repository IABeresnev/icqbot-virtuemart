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