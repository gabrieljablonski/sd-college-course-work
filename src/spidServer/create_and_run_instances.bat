go build spidServer.go
create_server_instances.py -d %1 -p %2
powershell -noexit server_mapper.py -d %1
run_servers.py -m %3
