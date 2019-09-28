import logging
from uuid import UUID

from request_handling.request_handler import RequestHandler

if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG,
                        format='%(asctime)s %(name)s: %(levelname)s >>> %(message)s',
                        datefmt='%y-%m-%d %H:%M:%S')

    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID client')
    parser.add_argument('--host', type=str, help='Server to connect to')
    parser.add_argument('-p', '--port', type=int, help='Server port')

    args = parser.parse_args()

    handler = RequestHandler(args.host, args.port)
    handler.connect()
    s = handler.get_spid_info(UUID("13e711ab-8021-4d2f-94e3-d20017fd3951"))


    handler.close_connection()