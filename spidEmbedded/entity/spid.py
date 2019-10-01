from uuid import UUID
from datetime import datetime
import json

from entity.spid_attributes import LockInfo, GlobalPosition
from utils import UUIDEncoder


class Spid:
    def __init__(
            self,
            id: str = "0"*32,
            battery_level: int = None,
            lock: LockInfo = None,
            location: GlobalPosition = None,
            last_updated: datetime = None,
            current_user_id: str = "0"*32):
        self.id: UUID = UUID(hex=id)
        self.battery_level = battery_level
        self.lock = lock
        self.location = location
        self.last_updated = last_updated
        self.current_user_id: UUID = UUID(hex=current_user_id)

    @classmethod
    def from_dict(cls, _dict):
        return cls(**_dict)

    @classmethod
    def from_json(cls, json_data):
        return cls(**json.loads(json_data))

    def to_json(self, indent=None):
        return json.dumps(self.__dict__, cls=UUIDEncoder, indent=indent)
