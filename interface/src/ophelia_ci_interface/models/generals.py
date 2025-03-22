from typing import Literal, Self

from pydantic import BaseModel

type MetadataTuple = tuple[tuple[str, str]]


class ModalItem(BaseModel):
    id: str
    label: str
    type: str
    autocomplete: str = 'off'
    options: list[str] = []


class Modal(BaseModel):
    title: str
    action: str
    method: Literal['POST', 'PUT', 'DELETE', 'GET']
    items: list[ModalItem] = []
    submit: str
    submit_id: str

    def format_action(self, **kwargs) -> Self:
        self.action = self.action.format(**kwargs)
        self.submit_id = self.submit_id.format(**kwargs)
        return self
