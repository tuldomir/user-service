#!/bin/bash
set -e

clickhouse client -n <<-EOSQL

CREATE DATABASE events;
CREATE TABLE userevent(

	event_type String,
	id UUID,
	email String,
	created_at DateTime

) ENGINE = MergeTree ORDER BY (event_type, created_at);


CREATE TABLE events.userevent_queue
(
        event_type String,
	id UUID,
	email String,
	created_at DateTime

) ENGINE = Kafka()
SETTINGS
	kafka_broker_list = '$KAFKA_HOST:$KAFKA_PORT',
	kafka_topic_list = 'useraddtopic',
	kafka_group_name = 'clickhouse',
	kafka_format = 'Protobuf',
	kafka_row_delimiter = '',
	kafka_schema = 'event:User',
	kafka_num_consumers = 1,
	kafka_thread_per_consumer = 0;


CREATE MATERIALIZED VIEW events.userevent_mv TO events.userevent AS \
SELECT * FROM events.userevent_queue;

EOSQL
