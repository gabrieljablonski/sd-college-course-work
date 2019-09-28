import json
from uuid import UUID

from utils import UUIDEncoder


class Spid:
    def __init__(self, id: str = "0"*32, battery_level: int = None, lock_state: str = None):
        self.id = UUID(hex=id)
        self.battery_level = battery_level
        self.lock_state = lock_state

    @classmethod
    def from_dict(cls, _dict):
        return cls(**_dict)

    @classmethod
    def from_json(cls, json_data):
        return cls(**json.loads(json_data))

    def to_json(self, indent=None):
        return json.dumps(self.__dict__, cls=UUIDEncoder, indent=indent)
