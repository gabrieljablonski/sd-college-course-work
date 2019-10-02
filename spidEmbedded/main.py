import logging
from uuid import UUID

from request_handling.request_handler import RequestHandler, Spid


def main(host, port):
    handler = RequestHandler(host, port)
    handler.connect()

    spid = Spid()
    valid = False
    while True:
        cmd = input('>> ').strip()
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

        elif cmd == 'load':
            fn = input('<< filename: ')
            try:
                with open(fn) as f:
                    uid = f.readline()
            except FileNotFoundError:
                print(f"-- file {fn} does not exist")
                continue
            try:
                uid = UUID(hex=uid)
            except ValueError as e:
                print(f"-- invalid id: `{e}`")
            else:
                spid.id = uid
                valid = True
                print(f"-- loaded id `{spid.id.hex}`")

        elif cmd == 'register':
            s = handler.register_spid()
            if s.id.int != 0:
                spid = s
                print(f"-- registered")
                valid = True
                continue
            print('-- failed to register')

        elif cmd == 'query':
            if not valid:
                uuid = input('<< id (uuid hex): ')
                try:
                    uuid = UUID(uuid)
                except ValueError as e:
                    print(f"-- invalid uuid: `{e}`")
                    continue
            else:
                uuid = spid.id
            s = handler.get_spid_info(uuid)
            if spid.id.int == 0:
                print('-- spid not found')
                continue
            spid = s
            valid = True
            print('-- spid found')

        elif cmd == 'save':
            if not valid:
                print('-- register or query first')
                continue
            fn = input('<< filename: ')
            try:
                with open(fn, 'w') as f:
                    f.write(spid.id.hex)
            except FileNotFoundError:
                print(f"-- file {fn} does not exist")
                continue
            print(f"-- wrote spid id to file {fn}")

        elif cmd == 'update location':
            if not valid:
                print('-- register or query first')
                continue
            lat = input('<< latitude (-90,90): ')
            lon = input('<< longitude (-180,180): ')
            try:
                lat = float(lat)
                lon = float(lon)
            except TypeError:
                print('-- invalid latitude or longitude values')
            else:
                if any((lat>90, lat<-90, lon>180, lon<-180)):
                    print('-- invalid latitude or longitude values')
                    continue
            old_lat = spid.location["latitude"]
            old_lon = spid.location["longitude"]
            spid.location["latitude"] = lat
            spid.location["longitude"] = lon
            s = handler.update_spid_location(spid)
            if s.id.int != 0:
                spid = s
                print('-- location updated')
                continue
            spid.location["latitude"] = old_lat
            spid.location["longitude"] = old_lon
            print(f"-- failed to update location")

        elif cmd == 'delete':
            if not valid:
                print('-- register or query first')
                continue
            s = handler.delete_spid(spid.id)
            if s.id.int != 0:
                spid = s
                print('-- spid deleted')
                continue
            print(f"-- failed to delete spid")

        else:
            available_commands = '\n\t-'.join((
                'view',
                'load',
                'register',
                'query',
                'save',
                'update location',
                'delete',
                'exit',
            ))
            print(f"-- Available commands:\n\t-"
                  f"{available_commands}")
    handler.close_connection()


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG,
                        format='%(asctime)s %(name)s: %(levelname)s >>> %(message)s',
                        datefmt='%y-%m-%d %H:%M:%S')

    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID client')
    parser.add_argument('--host', type=str, required=True, help='Server to connect to')
    parser.add_argument('-p', '--port', type=int, required=True, help='Server port')

    args = parser.parse_args()

    main(args.host, args.port)
