import os
import shutil

base_port = 45678
registrar_address = 'localhost', '43210'
base_dir = os.path.abspath('_server_instances')

number_of_servers = 9

if os.path.exists(base_dir):
    shutil.rmtree(base_dir)
os.mkdir(base_dir)

for i in range(number_of_servers):
    dst = os.path.join(base_dir, f"s{i}")
    if os.path.exists(dst):
        shutil.rmtree(dst)
    shutil.copytree(os.path.abspath('data'), os.path.join(dst, "data"))
    shutil.copy(os.path.abspath('spidServer.exe'), dst)

    with open(os.path.join(dst, "run.bat"), 'w') as f:
        f.write(f"powershell -noexit {os.path.join(dst, 'spidServer.exe')} {base_port+i} {registrar_address[0]} {registrar_address[1]}\n")
