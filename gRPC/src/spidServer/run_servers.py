import os
import subprocess

base_dir = os.path.abspath('_server_instances')
for i, dir in enumerate(os.listdir(base_dir)):
    subprocess.Popen([f"{base_dir}\\s{i}\\run.bat"], creationflags=subprocess.CREATE_NEW_CONSOLE)