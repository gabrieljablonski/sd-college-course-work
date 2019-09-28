from uuid import UUID, uuid4
import json

from utils import UUIDEncoder


class RequestType:
    GET_SPID_INFO = "GET SPID INFO"
    REGISTER_SPID = "REGISTER SPID"
    UPDATE_SPID_LOCATION = "UPDATE SPID LOCATION"
    DELETE_SPID = "DELETE SPID"
    REQUEST_LOCK_CHANGE = "REQUEST LOCK CHANGE"

    TIMEOUT = "TIMEOUT"
    INVALID = "INVALID"


class Request:
    def __init__(
            self,
            request_id: UUID = None,
            request_type: str = None,
            request_body: dict = None):
        self.id = request_id
        self.type = request_type
        self.body = request_body

    @classmethod
    def from_json(cls, json_data):
        return cls(**json.loads(json_data))

    @property
    def as_json(self):
        return json.dumps(self.__dict__, cls=UUIDEncoder)


class Response:
    def __init__(
            self,
            id: UUID = None,
            type: str = None,
            ok: bool = None,
            body: dict = None):
        self.id = id
        self.type = type
        self.ok = ok
        self.body = body

    @classmethod
    def from_json(cls, json_data):
        return cls(**json.loads(json_data))

    @property
    def as_json(self):
        return json.dumps(self.__dict__, cls=UUIDEncoder)
