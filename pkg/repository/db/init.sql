
CREATE TABLE orders (
                        order_uid VARCHAR PRIMARY KEY,
                        track_number VARCHAR,
                        entry VARCHAR,
                        locale VARCHAR,
                        internal_signature VARCHAR,
                        customer_id VARCHAR,
                        delivery_service VARCHAR,
                        shardkey VARCHAR,
                        sm_id INT,
                        date_created VARCHAR,
                        oof_shard VARCHAR
);

CREATE TABLE delivery (
                          order_uid VARCHAR REFERENCES orders(order_uid),
                          name VARCHAR,
                          phone VARCHAR,
                          zip VARCHAR,
                          city VARCHAR,
                          address VARCHAR,
                          region VARCHAR,
                          email VARCHAR
);

CREATE TABLE payment (
                         order_uid VARCHAR REFERENCES orders(order_uid),
                         transaction VARCHAR,
                         request_id VARCHAR,
                         currency VARCHAR,
                         provider VARCHAR,
                         amount INT,
                         payment_dt VARCHAR,
                         bank VARCHAR,
                         delivery_cost INT,
                         goods_total INT,
                         custom_fee INT
);

CREATE TABLE items (
                       chrt_id INT PRIMARY KEY,
                       order_uid VARCHAR REFERENCES orders(order_uid),
                       track_number VARCHAR,
                       price INT,
                       rid VARCHAR,
                       name VARCHAR,
                       sale INT,
                       size VARCHAR,
                       total_price INT,
                       nm_id INT,
                       brand VARCHAR,
                       status INT
);

DELETE FROM delivery WHERE order_uid = 'b563feb7b2b84b6test';
DELETE FROM payment WHERE order_uid = 'b563feb7b2b84b6test';
DELETE FROM items WHERE order_uid = 'b563feb7b2b84b6test';
DELETE FROM orders WHERE order_uid = 'b563feb7b2b84b6test';

drop table delivery, payment, items, orders;