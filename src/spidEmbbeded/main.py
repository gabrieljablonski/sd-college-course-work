import logging
from uuid import UUID

from request_handling.request_handler import RequestHandler, Spid


def main(host, port):
    handler = RequestHandler(host, port)
    handler.connect()

    spid = Spid()
    valid = False
    while True:
        cmd = input('>> ')
        if not cmd:
            continue
        if cmd == 'exit':
            break
        if cmd == 'view':
            if not valid:
                print('-- register or query first')
                continue
            print(f"-- spid:\n`{spid.to_json(4)}`")
            print(f"-- spid lock is {spid.lock['state']}")
            if spid.current_user_id.int != 0:
                print(f"-- spid is associated to user {spid.current_user_id.hex}")
            else:
                print('-- no user associated')
        if cmd == 'register':
            spid = handler.register_spid()
            valid = True
            print('-- registered')
        if cmd == 'query':
            uuid = input('<< id: ')
            spid = handler.get_spid_info(UUID(uuid))
            if spid.id.int == 0:
                print('-- spid not found')
                continue
            valid = True
            print('-- spid found')
        if cmd == 'update':
            if not valid:
                print('-- register or query first')
                continue
            lat = input('<< latitude: ')
            lon = input('<< longitude: ')
            try:
                lat = float(lat)
                lon = float(lon)
            except TypeError:
                print('-- invalid latitude or longitude values')
            else:
                if any((lat>90, lat<-90, lon>180, lon<-180)):
                    print('-- invalid latitude or longitude values')
                    continue
            spid.location["latitude"] = lat
            spid.location["longitude"] = lon
            spid = handler.update_spid_location(spid)
            print('-- location updated')
    handler.close_connection()


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG,
                        format='%(asctime)s %(name)s: %(levelname)s >>> %(message)s',
                        datefmt='%y-%m-%d %H:%M:%S')

    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID client')
    parser.add_argument('--host', type=str, help='Server to connect to')
    parser.add_argument('-p', '--port', type=int, help='Server port')

    args = parser.parse_args()

    main(args.host, args.port)
