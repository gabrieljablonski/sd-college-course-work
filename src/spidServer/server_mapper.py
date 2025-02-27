import socket
import json
import re


DEFAULT_PORT = 43210
DEFAULT_IP_MAP_PATH = 'ip_map.spd'


def main(port, base_delta, ip_map_path):
    if port is None:
        port = DEFAULT_PORT
    if ip_map_path is None:
        ip_map_path = ''
    number_of_servers = base_delta**2
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
    print(f"Bound to {addr[0]}:{addr[1]}")
    while True:
        print(f"listening...")
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
            try:
                uuid, server_port = recv.split()[2:]
            except IndexError as e:
                response = str(e)
            else:
                if len(ip_map) == number_of_servers and uuid not in ip_map:
                    # all server slots already filled
                    response = 'full'
                else:
                    ip_map[uuid] = f"{client_addr[0]}:{server_port}"
                    # assuming ordered dictionary
                    response = f"{list(ip_map).index(uuid)} {number_of_servers}"

                    if len(ip_map) and not ip_map_path:
                        with open(DEFAULT_IP_MAP_PATH, 'w') as f:
                            print(f"saving to file: {ip_map}")
                            json.dump(ip_map, f)
        
        elif 'REQUEST IP MAP' in recv:
            try:
                server_number = int(recv.split()[-1])
            except (ValueError, TypeError) as e:
                response = str(e)
            else:
                if len(ip_map) != number_of_servers:
                    response = '{}'
                else:
                    # assuming ordered dictionary
                    ip_list = list(ip_map.values())
                    # creates a map with the server ips in all 6 main cardinal directions
                    # from the server requesting the map.
                    #
                    # in all directions, the distance increases by a factor of 2 each time
                    # (first is 1 away, second is 2, third is 4, fourth is 8, ...).
                    #
                    # moving diagonally is also considered 1 unit of distance.
                    #
                    # the maximum size for a resulting map is bounded by:
                    # 6*floor(log2(n-1))
                    # which corresponds to a server exactly in the middle of the map
                    response = {str(server_number): ip_list[server_number]}
                    base_delta = int(round(number_of_servers**.5))
                    sX = server_number % base_delta
                    sY = server_number//base_delta
                    # north
                    tX, tY = sX, sY
                    b = 1
                    while 1:
                        tY = sY - b
                        if tY < 0:
                            break
                        tN = tY*base_delta + tX
                        response[str(tN)] = ip_list[tN]
                        b *= 2
                    # south
                    tX, tY = sX, sY
                    b = 1
                    while 1:
                        tY = sY + b
                        if tY >= base_delta:
                            break
                        tN = tY*base_delta + tX
                        response[str(tN)] = ip_list[tN]
                        b *= 2
                    # west
                    tX, tY = sX, sY
                    b = 1
                    while 1:
                        tX = sX - b
                        if tX < 0:
                            break
                        tN = tY*base_delta + tX
                        response[str(tN)] = ip_list[tN]
                        b *= 2
                    # east
                    tX, tY = sX, sY
                    b = 1
                    while 1:
                        tX = sX + b
                        if tX >= base_delta:
                            break
                        tN = tY*base_delta + tX
                        response[str(tN)] = ip_list[tN]
                        b *= 2
                    # northeast
                    tX, tY = sX, sY
                    b = 1
                    while 1:
                        tX = sX + b
                        tY = sY - b
                        if tX >= base_delta or tY < 0:
                            break
                        tN = tY*base_delta + tX
                        response[str(tN)] = ip_list[tN]
                        b *= 2
                    # southeast
                    tX, tY = sX, sY
                    b = 1
                    while 1:
                        tX = sX + b
                        tY = sY + b
                        if tX >= base_delta or tY >= base_delta:
                            break
                        tN = tY*base_delta + tX
                        response[str(tN)] = ip_list[tN]
                        b *= 2
                    # southwest
                    tX, tY = sX, sY
                    b = 1
                    while 1:
                        tX = sX - b
                        tY = sY + b
                        if tX < 0 or tY >= base_delta:
                            break
                        tN = tY*base_delta + tX
                        response[str(tN)] = ip_list[tN]
                        b *= 2
                    # northwest
                    tX, tY = sX, sY
                    b = 1
                    while 1:
                        tX = sX - b
                        tY = sY - b
                        if tX < 0 or tY < 0:
                            break
                        tN = tY*base_delta + tX
                        response[str(tN)] = ip_list[tN]
                        b *= 2
        print(f"sending {response}")
        response = str(response).replace("'",'"')
        conn.sendall(f"{response}\n".encode())
        conn.close()


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID server mapper')
    parser.add_argument('-d', '--base-delta', type=int, required=True, help='Base delta (number of map lines/columns)')
    parser.add_argument('-p', '--port', type=int, required=False, help='Mapper port')
    parser.add_argument('-t', '--ip-map', type=str, required=False, help='Path to already existing ip map')

    args = parser.parse_args()

    main(port=args.port, base_delta=args.base_delta, ip_map_path=args.ip_map)
