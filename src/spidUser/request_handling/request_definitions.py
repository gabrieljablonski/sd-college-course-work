from uuid import UUID, uuid4
import json

from utils import UUIDEncoder


class RequestType:
    GET_USER_INFO = "GET USER INFO"
    REGISTER_USER = "REGISTER USER"
    UPDATE_USER_LOCATION = "UPDATE USER LOCATION"
    DELETE_USER = "DELETE USER"

    REQUEST_ASSOCIATION = "REQUEST ASSOCIATION"
    REQUEST_DISSOCIATION = "REQUEST DISSOCIATION"
    REQUEST_SPID_INFO = "REQUEST SPID INFO"
    REQUEST_LOCK_CHANGE = "REQUEST LOCK CHANGE"

    TIMEOUT = "TIMEOUT"
    INVALID = "INVALID"


class Request:
    def __init__(
            self,
            id: UUID = None,
            type: str = None,
            body: dict = None):
        self.id = id
        self.type = type
        self.body = body

    @classmethod
    def from_json(cls, json_data):
        return cls(**json.loads(json_data))

    def to_json(self, indent=None):
        return json.dumps(self.__dict__, cls=UUIDEncoder, indent=indent)


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
