class LockInfo:
    def __init__(
            self,
            override: bool = None,
            pending: bool = None,
            state: str = None):
        self.override = override
        self.pending = pending
        self.state = state


class GlobalPosition:
    def __init__(
            self,
            latitude: float = None,
            longitude: float = None):
        self.latitude = latitude
        self.longitude = longitude
