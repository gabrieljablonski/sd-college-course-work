import os
import shutil

def main(base_delta=3, base_port=45678):
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
            # %1 and %2 -> mapper address and port
            f.write(f"powershell -noexit {os.path.join(dst, 'spidServer.exe')} {base_port+i} %1 %2\n")

if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID server instances')
    parser.add_argument('-d', '--delta', type=int, required=True, help='Base delta (number of map lines/columns)')
    parser.add_argument('-p', '--port', type=int, required=True, help='Base server port number')

    args = parser.parse_args()

    main(base_delta=args.delta, base_port=args.port)
