from uuid import UUID
from datetime import datetime
import json

from entity.user_attributes import GlobalPosition
from utils import UUIDEncoder


class User:
    def __init__(
            self,
            id: str = "0"*32,
            name: str = None,
            location: GlobalPosition = None,
            last_updated: datetime = None,
            current_spid_id: str = "0"*32):
        self.id: UUID = UUID(hex=id)
        self.name = name
        self.location = location
        self.last_updated = last_updated
        self.current_user_id: UUID = UUID(hex=current_spid_id)

    @classmethod
    def from_dict(cls, _dict):
        return cls(**_dict)

    @classmethod
    def from_json(cls, json_data):
        return cls(**json.loads(json_data))

    def to_json(self, indent=None):
        return json.dumps(self.__dict__, cls=UUIDEncoder, indent=indent)
