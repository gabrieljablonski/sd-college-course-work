import os
import subprocess

def main(mapper_address, mapper_port):
    base_dir = os.path.abspath('_server_instances')
    for i, dir in enumerate(os.listdir(base_dir)):
        subprocess.Popen([f"{base_dir}\\s{i}\\run.bat", mapper_address, str(mapper_port)], creationflags=subprocess.CREATE_NEW_CONSOLE)

if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID server instances')
    parser.add_argument('-a', '--mapper-address', type=str, required=True, help='Server mapper address')
    parser.add_argument('-p', '--mapper-port', type=int, required=True, help='Server mapper port')

    args = parser.parse_args()

    main(args.mapper_address, args.mapper_port)
