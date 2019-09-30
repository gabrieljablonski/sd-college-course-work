import logging
from uuid import UUID

from request_handling.request_handler import RequestHandler, User


def main(host, port):
    handler = RequestHandler(host, port)
    handler.connect()

    user = User()
    valid = False
    while True:
        cmd = input('>> ').strip()

        if cmd == 'exit':
            break

        if cmd == 'view user':
            if not valid:
                print('-- register or query first')
                continue
            print(f"-- user:\n`{user.to_json(4)}`")
            if user.current_spid_id.int != 0:
                print(f"-- user is associated to spid {user.current_spid_id.hex}")
            else:
                print('-- no spid associated')

        elif cmd == 'load user':
            fn = input('<< filename: ')
            with open(fn) as f:
                uid = f.readline()
            try:
                uid = UUID(hex=uid)
            except ValueError as e:
                print(f"-- invalid id: `{e}`")
            else:
                print(f"-- loaded id `{uid.hex}`. ready to query")
                user.id = uid
                valid = True

        elif cmd == 'register user':
            user_name = input('<< user name (ascii): ')
            user = handler.register_user(user_name)
            valid = True
            print(f"-- registered user `{user_name}`")

        elif cmd == 'query user':
            if not valid:
                uuid = input('<< id (uuid hex): ')
                try:
                    uuid = UUID(uuid)
                except ValueError as e:
                    print(f"-- invalid uuid: `{e}`")
                    continue
            else:
                uuid = user.id
            user = handler.get_user_info(uuid)
            if user.id.int == 0:
                print('-- user not found')
                continue
            valid = True
            print('-- user found')

        elif cmd == 'save user':
            if not valid:
                print('-- register or query first')
                continue
            fn = input('<< filename: ')
            with open(fn, 'w') as f:
                f.write(user.id.hex)
            print(f"-- wrote user id to file {fn}")

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
            user.location["latitude"] = lat
            user.location["longitude"] = lon
            user = handler.update_user_location(user)
            print('-- location updated')

        elif cmd == 'delete user':
            if not valid:
                print('-- register or query first')
                continue
            user = handler.delete_user(user.id)
            print('-- user deleted')

        elif cmd == 'associate spid':
            if user.current_spid_id.int != 0:
                print(f"-- user already associated to spid `{user.current_spid_id}`")
                continue
            option = ""
            while option != 'y' and option != 'n':
                option = input(f"<< load spid id from file? (Y/n): ").lower()
                if not option:
                    option = 'y'
            if option == 'y':
                fn = input('<< filename: ')
                with open(fn) as f:
                    uid = f.readline()
            else:
                uid = input('<< id (uuid hex): ')
            try:
                uid = UUID(hex=uid)
            except ValueError as e:
                print(f"-- invalid id: `{e}`")
                continue
            print(f"-- loaded id `{uid.hex}`. ready to query")
            user.current_spid_id = uid
            user = handler.request_association(user.id, user.current_spid_id)
            print(f"-- user associated to spid `{user.current_spid_id}`")

        elif cmd == 'save spid':
            if user.current_spid_id.int == 0:
                print('-- user not associated to any spids')
                continue
            fn = input('<< filename: ')
            with open(fn, 'w') as f:
                f.write(user.current_spid_id.hex)
            print(f"-- spid id saved to file {fn}")

        elif cmd == 'dissociate':
            if user.current_spid_id.int == 0:
                print('-- user not associated to any spids')
                continue
            user = handler.request_dissociation(user.id, user.current_spid_id)
            print(f"-- user dissociated from spid `{user.current_spid_id}`")

        elif cmd == 'query spid':
            if not valid:
                print('-- register or query user first')
                continue
            if user.current_spid_id.int == 0:
                print('-- user not associated to any spids')
                continue
            user.current_spid = handler.get_spid_info(user.id, user.current_spid_id)
            print(f"-- spid info:\n`{user.current_spid.to_json(4)}`")

        elif cmd == 'lock spid':
            if not valid:
                print('-- register or query user first')
                continue
            if user.current_spid_id.int == 0:
                print('-- user not associated to any spids')
                continue
            try:
                spid = handler.change_lock_state(user.id, user.current_spid_id, "locked")
            except Exception as e:
                print(f"-- failed to change lock state: `{e}`")
            else:
                user.current_spid = spid
                print('-- changed lock state to locked')

        elif cmd == 'unlock spid':
            if not valid:
                print('-- register or query user first')
                continue
            if user.current_spid_id.int == 0:
                print('-- user not associated to any spids')
                continue
            try:
                spid = handler.change_lock_state(user.id, user.current_spid_id, "unlocked")
            except Exception as e:
                print(f"-- failed to change lock state: `{e}`")
            else:
                user.current_spid = spid
                print('-- changed lock state to unlocked')

        else:
            available_commands = '\n\t-'.join((
                'view user',
                'load user',
                'register user',
                'query user',
                'save user',
                'update location',
                'delete user',
                'associate spid',
                'save spid',
                'dissociate',
                'query spid',
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
