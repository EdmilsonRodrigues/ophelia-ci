from pydantic import BaseModel


class ModalItem(BaseModel):
    id: str
    label: str
    type: str
    autocomplete: str = 'off'
    options: list[str] = []


class Modal(BaseModel):
    title: str
    items: list[ModalItem] = []
    submit: str
    submit_id: str
