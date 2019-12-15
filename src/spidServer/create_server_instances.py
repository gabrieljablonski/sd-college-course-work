import os
import shutil


def main(base_delta=3, base_port=45678, nodes_per_cluster=3):
    base_dir = os.path.abspath('_server_instances')
    number_of_servers = base_delta**2

    if os.path.exists(base_dir):
        shutil.rmtree(base_dir)
    os.mkdir(base_dir)

    for i in range(number_of_servers):
        s = os.path.join(base_dir, f"s{i}")
        os.mkdir(s)
        for j in range(nodes_per_cluster):
            n = os.path.join(s, f"c{j}")
            os.mkdir(n)
            shutil.copy(os.path.abspath('spidServer.exe'), n)
            with open(os.path.join(n, "run.bat"), 'w') as f:
                # %1 -> mapper address
                # %2 -> cluster endpoints
                exe_path = os.path.join(n, 'spidServer.exe')
                to_write = (
                    ":: %1 -> server mapper address",
                    ":: %2 -> cluster endpoints",
                    f"powershell -noexit {exe_path} {base_port+nodes_per_cluster*i+j} %1 %2",
                )
                f.write('\n'.join(to_write)+'\n')


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID server instances')
    parser.add_argument('-d', '--delta', type=int, required=True, help='Base delta (number of map lines/columns)')
    parser.add_argument('-p', '--port', type=int, required=True, help='Base server port number')
    parser.add_argument('-n', '--nodes-per-cluster', type=int, required=True, help='Number of nodes per cluster')

    args = parser.parse_args()

    main(base_delta=args.delta, base_port=args.port, nodes_per_cluster=args.nodes_per_cluster)
