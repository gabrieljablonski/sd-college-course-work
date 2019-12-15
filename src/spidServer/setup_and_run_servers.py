import os
import subprocess
from shutil import rmtree


PORTS_PER_CLUSTER = 100
BASE_PATH = r'D:\GitReps\sd-college-course-work\src\spidServer'
BOOT_ETCD = f"{BASE_PATH}\\boot_etcd.bat"
CREATE_AND_RUN = f"{BASE_PATH}\\create_and_run_instances.bat"

ETCD_FILES = r'D:\GitReps\sd-college-course-work\etcd\clusters'
ETCD = r'D:\GitReps\sd-college-course-work\etcd\etcd.exe'


def main(delta, servers_base_port, server_mapper_address, clusters_base_port, nodes_per_cluster):
    cluster_names = []
    node_names = []
    endpoints = []
    endpoint_strings = []

    cluster_port = clusters_base_port
    for i in range(delta*delta):
        cluster_names.append(f"etcd-cluster-{i}")
        node_names.append([])
        endpoints.append([])

        endpoint_strings.append("")
        node_port = cluster_port
        for j in range(nodes_per_cluster):
            node_names[i].append(f"infra_{j}")
            endpoints[i].append(f"http://localhost:{node_port}")
            node_port += 1

            endpoint_strings[i] += f"{node_names[i][j]}={endpoints[i][j]}"
            endpoint_strings[i] += ',' if j != nodes_per_cluster-1 else ''
        cluster_port += PORTS_PER_CLUSTER
    base_client_port = 9999

    if os.path.exists(ETCD_FILES):
        rmtree(ETCD_FILES)
    os.mkdir(ETCD_FILES)
    k = 0
    for i, (c, ns, es, ess) in enumerate(zip(cluster_names, node_names, endpoints, endpoint_strings)):
        os.chdir(ETCD_FILES)
        c_path = f"cluster_{i}"
        os.mkdir(c_path)
        os.chdir(c_path)
        for n, e in zip(ns, es):
            subprocess.Popen([BOOT_ETCD, n, e, c, f'"{ess}"', f"http://localhost:{base_client_port+k}"],
                             creationflags=subprocess.CREATE_NEW_CONSOLE)
            k += 1
    os.chdir(BASE_PATH)

    subprocess.call(['go', 'build', 'spidServer.go'])
    subprocess.call(['python', 'create_server_instances.py',
                     '-d', str(delta),
                     '-p', str(servers_base_port),
                     '-n', str(nodes_per_cluster)])
    subprocess.Popen(['server_mapper.bat',
                      str(server_mapper_address.split(':')[1]),
                      str(delta),
                      str(nodes_per_cluster)],
                     creationflags=subprocess.CREATE_NEW_CONSOLE)

    base_dir = os.path.abspath('_server_instances')
    for i, (_, es) in enumerate(zip(os.listdir(base_dir), endpoints)):
        for j, e in enumerate(es):
            subprocess.Popen([f"{base_dir}\\s{i}\\c{j}\\run.bat", server_mapper_address, e],
                             creationflags=subprocess.CREATE_NEW_CONSOLE)


if __name__ == '__main__':
    import argparse

    parser = argparse.ArgumentParser(description='Setup SPID server instances')
    parser.add_argument('-d', '--delta', type=int, required=True, help='Base delta')
    parser.add_argument('-p', '--servers-base-port', type=int, required=True, help='Servers base port')
    parser.add_argument('-m', '--server-mapper-address', type=str, required=True, help='Address for server mapper')
    parser.add_argument('-c', '--clusters-base-port', type=int, required=True, help='Clusters base port')
    parser.add_argument('-n', '--nodes-per-cluster', type=int, required=True, help='Number of nodes per cluster')

    args = parser.parse_args()

    main(
        args.delta,
        args.servers_base_port,
        args.server_mapper_address,
        args.clusters_base_port,
        args.nodes_per_cluster,
    )
