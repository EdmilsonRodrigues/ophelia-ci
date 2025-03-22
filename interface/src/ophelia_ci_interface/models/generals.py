import json
import logging
from typing import Literal, Self

from pydantic import BaseModel

type MetadataTuple = tuple[tuple[str, str]]


class ModalItem(BaseModel):
    """
    Modal Item Format

    Attributes:
        id: The id of the modal item
        label: The label of the modal item
        type: The type of the modal item
        autocomplete: The autocomplete of the modal item
        options: The options of the modal item
    """

    id: str
    label: str
    type: str
    autocomplete: str = 'off'
    options: list[str] = []


class Modal(BaseModel):
    """
    Modal Format

    Attributes:
        title: The title of the modal
        action: The action of the modal
        method: The method of the modal
        items: The items of the modal
        submit: The submit of the modal
    """

    title: str
    action: str
    method: Literal['POST', 'PUT', 'DELETE', 'GET']
    items: list[ModalItem] = []
    submit: str

    def format_action(self, **kwargs) -> Self:
        self.action = self.action.format(**kwargs)
        return self


class OpheliaException(Exception):
    """
    Ophelia Exception
    """

    pass


def log_formatted(
    event: str,
    *,
    logging_level: int = logging.INFO,
    sensitive: bool = False,
    **kwargs,
) -> None:
    logging.log(
        logging_level,
        json.dumps({
            'event': event,
            **kwargs,
            'sensitive': sensitive,
        }),
    )
