import logging
from uuid import uuid4, UUID
from json import JSONDecodeError

from tcp_client import TCPClient
from request_handling.request_definitions import RequestType as RT, Request, Response
from entity.user import User
from entity.spid import Spid


class RequestHandler:
    END_CONNECTION = 'END CONNECTION'
    END_CONNECTION_ACK = 'ENDING CONNECTION'

    def __init__(self, host, port):
        self._con = TCPClient(host, port)

    def connect(self, try_forever=False):
        self._con.connect(try_forever=try_forever)

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
        while True:
            try:
                self._con.send(request.to_json())
                return
            except ConnectionResetError as e:
                logging.error(e)
                logging.error('Lost connection to host, retrying...')
                self._con.connected = False
                self.connect(try_forever=True)

    def _get_response(self):
        try:
            return Response.from_json(self._con.receive())
        except (ConnectionResetError, JSONDecodeError) as e:
            logging.error('Lost connection to host, retrying...')
            self._con.connected = False
            self.connect(try_forever=True)
            return Response(body={})

    def get_user_info(self, uuid: UUID) -> User:
        request = Request(
            id=uuid4(),
            type=RT.GET_USER_INFO,
            body={"user_id": uuid}
        )
        logging.info(f"Querying user with id {uuid.hex}...")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            u = User.from_dict(response.body.get('user'))
            logging.info(f"User found.")
            return u
        logging.error(f"Failed to query user: `{response.body.get('message')}`")
        return User()

    def register_user(self, user_name: str) -> User:
        request = Request(
            id=uuid4(),
            type=RT.REGISTER_USER,
            body={"user_name": user_name}
        )
        logging.info(f"Registering new user with name `{user_name}`...")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            u = User.from_dict(response.body.get('user'))
            logging.info(f"Registered user.")
            return u
        logging.error(f"Failed to register user: `{response.body.get('message')}`")
        return User()

    def update_user_location(self, user: User) -> User:
        request = Request(
            id=uuid4(),
            type=RT.UPDATE_USER_LOCATION,
            body={
                "user_id": user.id,
                "location": user.location,
            }
        )
        logging.info(f"Updating user {user.id} location to {user.location}...")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            u = User.from_dict(response.body.get('user'))
            logging.info(f"Updated user.")
            return u
        logging.error(f"Failed to update user: `{response.body.get('message')}`")
        return User()

    def delete_user(self, uid: UUID) -> User:
        request = Request(
            id=uuid4(),
            type=RT.DELETE_USER,
            body={"user_id": uid}
        )
        logging.info(f"Deleting user with id {uid}...")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            u = User.from_dict(response.body.get('user'))
            logging.info(f"Deleted user.")
            return u
        logging.error(f"Failed to delete user: `{response.body.get('message')}`")
        return User()

    def request_association(self, u_uuid: UUID, s_uuid: UUID) -> User:
        request = Request(
            id=uuid4(),
            type=RT.REQUEST_ASSOCIATION,
            body={"user_id": u_uuid, "spid_id": s_uuid}
        )
        logging.info(f"Associating user {u_uuid} to spid {s_uuid}...")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            u = User.from_dict(response.body.get('user'))
            logging.info(f"Associated.")
            return u
        logging.error(f"Failed to associate to spid: `{response.body.get('message')}`")
        return User()

    def request_dissociation(self, u_uuid: UUID, s_uuid: UUID) -> User:
        request = Request(
            id=uuid4(),
            type=RT.REQUEST_DISSOCIATION,
            body={"user_id": u_uuid, "spid_id": s_uuid}
        )
        logging.info(f"Dissociating user {u_uuid} from spid {s_uuid}...")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            u = User.from_dict(response.body.get('user'))
            logging.info('Dissociated.')
            return u
        logging.error(f"Failed to dissociate from spid: `{response.body.get('message')}`")
        return User()

    def get_spid_info(self, u_uuid: UUID, s_uuid: UUID) -> Spid:
        request = Request(
            id=uuid4(),
            type=RT.REQUEST_SPID_INFO,
            body={"user_id": u_uuid, "spid_id": s_uuid}
        )
        logging.info(f"Requesting info for spid {s_uuid}...")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            s = Spid.from_dict(response.body.get('spid'))
            logging.info(f"Spid info acquired.")
            return s
        logging.error(f"Failed to get info for spid: `{response.body.get('message')}`")
        return Spid()

    def change_lock_state(self, u_uuid: UUID, s_uuid: UUID, lock_state: str) -> Spid:
        request = Request(
            id=uuid4(),
            type=RT.REQUEST_LOCK_CHANGE,
            body={"user_id": u_uuid, "spid_id": s_uuid, "lock_state": lock_state}
        )
        logging.info(f"Changing lock state from spid {s_uuid} to {lock_state}")
        self._make_request(request)
        response = self._get_response()
        if response.ok:
            s = Spid.from_dict(response.body.get('spid'))
            logging.info('Lock state changed.')
            return s
        msg = f"Failed to change lock state: `{response.body.get('message')}`"
        logging.error(msg)
        raise Exception(msg)
