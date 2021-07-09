env = "development"  # production
log_level = "debug"
apiserver_port = ":12000"
send_to_chat = false
logstash_url = "elk.b2bpolis.ru:5000"

[storage]
backend = "consul"

[consul]
address = ""
token = ""
