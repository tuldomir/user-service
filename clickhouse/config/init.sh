#!/bin/bash
set -e

clickhouse client -n <<-EOSQL

CREATE DATABASE IF NOT EXISTS events;
CREATE TABLE IF NOT EXISTS userevent(

	event_type String,
	uid String,
	created_at DateTime64

) ENGINE = MergeTree ORDER BY (event_type, created_at);

SET date_time_input_format='best_effort';


CREATE TABLE IF NOT EXISTS events.userevent_queue
(
    event_type String,
	uid String,
	created_at DateTime64

) ENGINE = Kafka()
SETTINGS
	kafka_broker_list = '$KAFKA_HOST:$KAFKA_PORT',
	kafka_topic_list = 'useraddtopic',
	kafka_group_name = 'clickhouse',
	kafka_format = 'JSONEachRow',
    input_format_import_nested_json = 1,
    date_time_input_format='best_effort',
	kafka_row_delimiter = '',
	kafka_num_consumers = 1,
	kafka_thread_per_consumer = 0;


CREATE MATERIALIZED VIEW IF NOT EXISTS events.userevent_mv TO default.userevent AS \
SELECT * FROM events.userevent_queue;

EOSQL
