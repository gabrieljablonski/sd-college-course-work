import logging
from uuid import uuid4, UUID

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
        self._con.send(request.to_json())

    def _get_response(self):
        return Response.from_json(self._con.receive())

    def get_spid_info(self, uuid: UUID) -> Spid:
        request = Request(
            id=uuid4(),
            type=RT.GET_SPID_INFO,
            body={"spid_id": uuid}
        )
        logging.info(f"Querying spid with id {uuid.hex}...")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            s = Spid.from_dict(response.body.get("spid"))
            logging.info(f"Found spid: `{s.to_json()}`")
            return s
        logging.error(f"Failed to query spid: `{response.body.get('message')}`")
        return Spid()

    def register_spid(self) -> Spid:
        request = Request(
            id=uuid4(),
            type=RT.REGISTER_SPID,
        )
        logging.info('Registering new spid...')
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            s = Spid.from_dict(response.body.get("spid"))
            logging.info(f"Registered spid: `{s.to_json()}`")
            return s
        logging.error(f"Failed to register spid: `{response.body.get('message')}`")
        return Spid()

    def update_spid_location(self, spid: Spid):
        request = Request(
            id=uuid4(),
            type=RT.UPDATE_SPID_LOCATION,
            body={
                "spid_id": spid.id,
                "location": spid.location,
            }
        )
        logging.info(f"Updating spid {spid.id} location to {spid.location}...")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            s = Spid.from_dict(response.body.get("spid"))
            logging.info(f"Updated spid: `{s.to_json()}`")
            return s
        logging.error(f"Failed to update spid: `{response.body.get('message')}`")
        return Spid()

    def delete_spid(self, uid: UUID):
        request = Request(
            id=uuid4(),
            type=RT.UPDATE_SPID_LOCATION,
            body={"spid_id": uid}
        )
        logging.info(f"Deleting spid with id {uid}...")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            s = Spid.from_dict(response.body.get("spid"))
            logging.info(f"Deleted spid: `{s.to_json()}`")
            return s
        logging.error(f"Failed to delete spid: `{response.body.get('message')}`")
        return Spid()
