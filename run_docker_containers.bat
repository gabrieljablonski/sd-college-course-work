set /A delta = %1%
set /A n = delta * delta - 1

for /l %%i in (0, 1, %n%) do (
    start powershell -noexit ssh -i C:\Users\gabri\Documents\scripts\jablonski-vm-private-key jablonski@localhost -p 2222 -t "sudo docker start -ia spid_grpc_server%%i"
)
