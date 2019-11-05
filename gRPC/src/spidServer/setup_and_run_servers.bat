:: %1 -> base delta
:: %2 -> servers base port
:: %3 -> server mapper address
:: %4 -> server mapper port

go build spidServer.go
create_server_instances.py -d %1 -p %2
start server_mapper.py -d %1
run_servers.py -a %3 -p %4
