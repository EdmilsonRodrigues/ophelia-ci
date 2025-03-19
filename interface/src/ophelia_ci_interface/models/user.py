from functools import cache

from ophelia_ci_interface.config import Settings
from ophelia_ci_interface.services.gRPC_service import UserService
from pydantic import UUID4, BaseModel


class CreateUserRequest(BaseModel):
    user_username: str
    user_private_key: str


class UpdateUserRequest(BaseModel):
    id: UUID4
    user_username: str
    user_private_key: str


class User(BaseModel):
    model_config = {'arbitrary_types_allowed': True}

    id: UUID4
    username: str

    @staticmethod
    @cache
    def get_service() -> UserService:
        return UserService(str(Settings().GRPC_SERVER))

    @classmethod
    def create(
        cls,
        username: str,
        private_key: str,
        metadata: tuple[tuple[str, str]],
    ):
        response = cls.get_service().create_user(
            username, private_key, metadata=metadata
        )
        return cls(id=response.id, username=response.username)

    @classmethod
    def update(
        cls,
        id: str,
        username: str,
        private_key: str,
        metadata: tuple[tuple[str, str]],
    ):
        response = cls.get_service().update_user(
            id, username, private_key, metadata=metadata
        )
        return cls(id=response.id, username=response.username)

    @classmethod
    def get(cls, id: str, metadata: tuple[tuple[str, str]]):
        response = cls.get_service().get_user(id, metadata=metadata)
        return cls(id=response.id, username=response.username)

    @classmethod
    def get_by_username(cls, username: str, metadata: tuple[tuple[str, str]]):
        response = cls.get_service().get_user_by_username(
            username, metadata=metadata
        )
        return cls(id=response.id, username=response.username)

    @classmethod
    def delete(cls, id: str, metadata: tuple[tuple[str, str]]):
        cls.get_service().delete_user(id, metadata=metadata)

    @classmethod
    def get_all(cls, metadata: tuple[tuple[str, str]]):
        response_list = cls.get_service().get_users(metadata=metadata)
        return [
            cls(id=response.id, username=response.username)
            for response in response_list.users
        ]
