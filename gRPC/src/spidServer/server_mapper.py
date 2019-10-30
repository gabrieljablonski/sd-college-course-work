import socket
import json
import re


DEFAULT_PORT = 45678
DEFAULT_IP_TABLE_PATH = 'ip_table.spd'


def main(port, number_of_servers, ip_table_path):
    if port is None:
        port = DEFAULT_PORT
    if ip_table_path is None:
        ip_table_path = ''
    ip_table = {}
    if ip_table_path:
        with open(ip_table_path) as f:
            ip_table = json.load(f)
        if len(ip_table) != number_of_servers:
            print(f"Existing table does not contain {number_of_servers} "
                  f"server{'' if len(ip_table) == 1 else 's'}.")
            ip_table = {}
            ip_table_path = ''
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

            if len(ip_table) == number_of_servers and uuid not in ip_table:
                # all server slots already filled
                response = 'full'

            ip_table[uuid] = f"{client_addr[0]}:{server_port}"

            if len(ip_table) and not ip_table_path:
                with open(DEFAULT_IP_TABLE_PATH, 'w') as f:
                    print(f"saving to file: {ip_table}")
                    json.dump(ip_table, f)
            # assuming ordered dictionary
            response = list(ip_table).index(uuid)
        
        elif recv == 'REQUEST IP TABLE':
            if len(ip_table) != number_of_servers:
                response = '[]'
            else:
                response = [addr for addr in ip_table.values()]
        print(f"sending {response}")
        response = str(response).replace("'",'"')
        conn.sendall(f"{response}\n".encode())
        conn.close()


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID client')
    parser.add_argument('-p', '--port', type=int, required=False, help='Mapper port')
    parser.add_argument('-n', '--number-of-servers', type=int, required=True, help='Number of servers to expect')
    parser.add_argument('-t', '--ip-table', type=str, required=False, help='Path to already existing ip table')

    args = parser.parse_args()

    main(port=args.port, number_of_servers=args.number_of_servers, ip_table_path=args.ip_table)