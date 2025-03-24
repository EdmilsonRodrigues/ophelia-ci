from typing import Annotated, Self

from ophelia_ci_interface.models.generals import MetadataTuple
from ophelia_ci_interface.services.gRPC_service import UserService
from pydantic import UUID4, BaseModel, Field


class User(BaseModel):
    """
    Model that represents a user.

    Attributes:
        id (UUID4): The ID of the user.
        username (str): The username of the user.
    """

    model_config = {'arbitrary_types_allowed': True}

    id: Annotated[UUID4, Field(title='ID', description='The ID of the user.')]
    username: Annotated[
        str, Field(title='Username', description='The username of the user.')
    ]

    @classmethod
    def create(
        cls,
        user_service: UserService,
        username: str,
        public_key: str,
        metadata: MetadataTuple,
    ) -> Self:
        """
        Create a new user in the database.

        :param user_service: the user service to use
        :param username: the username of the user
        :param public_key: the public key of the user
        :param metadata: the metadata of the request
        :return: the newly created user
        """
        response = user_service.create_user(
            username, public_key, metadata=metadata
        )
        return cls(id=response.id, username=response.username)

    @classmethod
    def update(
        cls,
        user_service: UserService,
        id: str,
        username: str,
        public_key: str,
        metadata: MetadataTuple,
    ) -> Self:
        """
        Update an existing user in the database.

        :param user_service: the user service to use
        :param id: the id of the user
        :param username: the username of the user
        :param public_key: the public key of the user
        :param metadata: the metadata of the request
        :return: the updated user
        """
        response = user_service.update_user(
            id, username, public_key, metadata=metadata
        )
        return cls(id=response.id, username=response.username)

    @classmethod
    def get(
        cls, user_service: UserService, id: str, metadata: MetadataTuple
    ) -> Self:
        """
        Get an existing user from the database.

        :param user_service: the user service to use
        :param id: the id of the user
        :param metadata: the metadata of the request
        :return: the user
        """
        response = user_service.get_user(id, metadata=metadata)
        return cls(id=response.id, username=response.username)

    @classmethod
    def get_by_username(
        cls, user_service: UserService, username: str, metadata: MetadataTuple
    ) -> Self:
        """
        Get an existing user from the database by its username.

        :param user_service: the user service to use
        :param username: the username of the user
        :param metadata: the metadata of the request
        :return: the user
        """
        response = user_service.get_user_by_username(
            username, metadata=metadata
        )
        return cls(id=response.id, username=response.username)

    @classmethod
    def delete(
        cls, user_service: UserService, id: str, metadata: MetadataTuple
    ) -> None:
        """
        Delete an existing user from the database.

        :param user_service: the user service to use
        :param id: the id of the user
        :param metadata: the metadata of the request
        :return: None
        """
        user_service.delete_user(id, metadata=metadata)

    @classmethod
    def get_all(
        cls, user_service: UserService, metadata: MetadataTuple
    ) -> list[Self]:
        """
        Get all existing users from the database.

        :param user_service: the user service to use
        :param metadata: the metadata of the request
        :return: the list of users
        """
        response_list = user_service.get_users(metadata=metadata)
        return [
            cls(id=response.id, username=response.username)
            for response in response_list.users
        ]
