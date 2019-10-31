import socket
import json
import re


DEFAULT_PORT = 43210
DEFAULT_IP_MAP_PATH = 'ip_map.spd'


def main(port, number_of_servers, ip_map_path):
    if port is None:
        port = DEFAULT_PORT
    if ip_map_path is None:
        ip_map_path = ''
    ip_map = {}
    if ip_map_path:
        with open(ip_map_path) as f:
            ip_map = json.load(f)
        if len(ip_map) != number_of_servers:
            print(f"Existing map does not contain {number_of_servers} "
                  f"server{'' if len(ip_map) == 1 else 's'}.")
            ip_map = {}
            ip_map_path = ''
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(1)
    addr = 'localhost', port

    s.bind(addr)
    s.listen(number_of_servers)
    while True:
        print('listening...')
        try:
            conn, client_addr = s.accept()
        except socket.timeout:
            continue
        print(f"serving {client_addr[0]}:{client_addr[1]}")

        recv = conn.recv(4096).decode().strip()
        print(recv)
        response = ''
        if 'REGISTER SERVER' in recv:
            # REGISTER SERVER <uuid> <port>
            uuid, server_port = recv.split()[2:]

            if len(ip_map) == number_of_servers and uuid not in ip_map:
                # all server slots already filled
                response = 'full'

            ip_map[uuid] = f"{client_addr[0]}:{server_port}"

            if len(ip_map) and not ip_map_path:
                with open(DEFAULT_IP_MAP_PATH, 'w') as f:
                    print(f"saving to file: {ip_map}")
                    json.dump(ip_map, f)
            # assuming ordered dictionary
            response = f"{list(ip_map).index(uuid)} {number_of_servers}"
        
        elif 'REQUEST IP MAP' in recv:
            if len(ip_map) != number_of_servers:
                response = '{}'
            else:
                response = {str(i):addr for i, addr in enumerate(ip_map.values())}
        print(f"sending {response}")
        response = str(response).replace("'",'"')
        conn.sendall(f"{response}\n".encode())
        conn.close()


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID client')
    parser.add_argument('-p', '--port', type=int, required=False, help='Mapper port')
    parser.add_argument('-n', '--number-of-servers', type=int, required=True, help='Number of servers to expect')
    parser.add_argument('-t', '--ip-map', type=str, required=False, help='Path to already existing ip map')

    args = parser.parse_args()

    main(port=args.port, number_of_servers=args.number_of_servers, ip_map_path=args.ip_map)