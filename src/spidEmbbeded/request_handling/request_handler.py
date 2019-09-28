import logging
from uuid import uuid4

from tcp_client import TCPClient
from request_handling.request_definitions import RequestType as RT, Request, Response
from entity.spid import Spid


class RequestHandler:
    END_CONNECTION = 'END CONNECTION'
    END_CONNECTION_ACK = 'ENDING CONNECTION'

    def __init__(self, host, port):
        self._con = TCPClient(host, port)

    def connect(self):
        self._con.connect()

    def close_connection(self):
        logging.info('Ending connection...')
        self._con.send(self.END_CONNECTION)
        response = self._con.receive()
        if response == self.END_CONNECTION_ACK:
            self._con.close()
            logging.info('Connection ended.')
        else:
            logging.critical('Failed to end connection.')
            
    def _make_request(self, request: Request):
        self._con.send(request.as_json)

    def _get_response(self):
        return Response.from_json(self._con.receive())

    def register_spid(self) -> Spid:
        request = Request(
            request_id=uuid4(),
            request_type=RT.REGISTER_SPID,
        )
        logging.info('Registering new spid...')
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            s = Spid.from_dict(response.body.get("spid"))
            logging.info(f"Registered spid: `{s.as_json}`")
            return s
