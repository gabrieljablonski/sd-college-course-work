import os
import shutil

def main(base_delta=3, base_port=45678, mapper_address='localhost', mapper_port=43210):
    base_dir = os.path.abspath('_server_instances')
    number_of_servers = base_delta**2

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
            f.write(f"powershell -noexit {os.path.join(dst, 'spidServer.exe')} {base_port+i} {mapper_address} {mapper_port}\n")

if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID server instances')
    parser.add_argument('-d', '--delta', type=int, required=True, help='Base delta (number of map lines/columns)')
    parser.add_argument('-p', '--port', type=int, required=True, help='Base server port number')
    parser.add_argument('-a', '--mapper-address', type=str, required=True, help='Server mapper address')
    parser.add_argument('-m', '--mapper-port', type=int, required=True, help='Server mapper port')

    args = parser.parse_args()

    main(base_delta=args.delta, base_port=args.port, mapper_address=args.mapper_address, mapper_port=args.mapper_port)
