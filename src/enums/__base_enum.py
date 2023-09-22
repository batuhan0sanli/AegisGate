from enum import Enum
from typing import Optional, List


class BaseEnum(Enum):
    @classmethod
    def get_by_name(cls, name):
        name = name.upper()
        for mode in cls:
            if mode.name == name:
                return mode
        raise ValueError(f"Invalid {cls.__name__} name: {name}")

    @classmethod
    def list_values(cls, exclude: Optional[List[Enum]] = None):
        exclude = exclude or []
        return [enum.value for enum in cls if enum.name not in exclude]
