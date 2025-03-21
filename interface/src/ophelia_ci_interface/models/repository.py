from datetime import datetime
from typing import Annotated, Self

from google.protobuf.timestamp_pb2 import Timestamp
from ophelia_ci_interface.config import Settings
from ophelia_ci_interface.models.generals import MetadataTuple
from ophelia_ci_interface.services.gRPC_service import RepositoryService
from pydantic import UUID4, BaseModel, Field, model_validator


class CreateRepositoryRequest(BaseModel):
    """
    Model that represents the request to create a new repository.

    Attributes:
        repository_name (str): The name of the repository.
        repository_description (str): The description of the repository.
        repository_gitignore (str): The main language of the repository, to be
            used for the base gitignore.
    """

    repository_name: Annotated[
        str, Field(title='Name', description='The name of the repository.')
    ]
    repository_description: Annotated[
        str,
        Field(
            title='Description',
            description='The description of the repository.',
        ),
    ]
    repository_gitignore: Annotated[
        str,
        Field(
            title='Gitignore',
            description='The main language of the repository, '
            'to be used for the base gitignore.',
        ),
    ]


class UpdateRepositoryRequest(BaseModel):
    """
    Model that represents the request to update an existing repository.

    Attributes:
        id (UUID4): The ID of the repository.
        repository_name (str): The name of the repository.
        repository_description (str): The description of the repository.
    """

    id: Annotated[
        UUID4, Field(title='ID', description='The ID of the repository.')
    ]
    repository_name: Annotated[
        str, Field(title='Name', description='The name of the repository.')
    ]
    repository_description: Annotated[
        str,
        Field(
            title='Description',
            description='The description of the repository.',
        ),
    ]


class Repository(BaseModel):
    """
    Model that represents a repository.

    Attributes:
        id (UUID4): The ID of the repository.
        name (str): The name of the repository.
        description (str): The description of the repository.
        last_updated (datetime): The last updated time of the repository.
        truncated_description (str): The truncated description of the
            repository.
        clone_url (str): The clone URL of the repository.
    """

    model_config = {'arbitrary_types_allowed': True}

    id: UUID4
    name: str
    description: str
    last_updated: datetime
    _clone_url: str | None = None

    @model_validator(mode='before')
    def parse_timestamp(cls, value):
        value['last_updated'] = cls.convert_timestamp_to_datetime(
            value['last_updated']
        )
        return value

    @property
    def truncated_description(self) -> str:
        """
        Get the truncated description of the repository.

        :return: the truncated description
        """
        if len(self.description) > 100:
            return self.description[:100] + '...'
        return self.description

    def get_clone_url(self, settings: Settings) -> None:
        """
        Generate the clone URL of the repository.

        This method generates the clone URL by combining the hostname of the
        gRPC server with the repository name. The generated URL is of the form
        `git@<hostname>:<repository_name>.git`.

        :param settings: the settings to use
        """
        self._clone_url = (
            f'git@{str(settings.GRPC_SERVER).rstrip("/")}:{self.name}.git'
        )

    @property
    def clone_url(self) -> str:
        """
        Get the clone URL of the repository.

        :return: the clone URL
        """
        if self._clone_url is None:
            raise ValueError(
                'Clone URL not generated. Call get_clone_url() first.'
            )
        return self._clone_url

    @staticmethod
    def convert_timestamp_to_datetime(timestamp: Timestamp) -> datetime:
        """
        Convert a google.protobuf.timestamp_pb2.Timestamp to a datetime object.

        :param timestamp: the Timestamp to convert
        :return: the converted datetime object
        """
        return timestamp.ToDatetime()

    @classmethod
    def create(
        cls,
        repository_service: RepositoryService,
        name: str,
        description: str,
        gitignore: str,
        metadata: MetadataTuple,
    ) -> Self:
        """
        Create a new repository in the database.

        :param repository_service: the repository service to use
        :param name: the name of the repository
        :param description: the description of the repository
        :param gitignore: the gitignore of the repository
        :param metadata: the metadata of the request
        :return: the newly created repository
        """
        response = repository_service.create_repository(
            name, description, gitignore, metadata=metadata
        )
        return cls(
            id=response.id,
            name=response.name,
            description=response.description,
            last_updated=response.last_update,
        )

    @classmethod
    def update(
        cls,
        repository_service: RepositoryService,
        id: str,
        name: str,
        description: str,
        metadata: MetadataTuple,
    ) -> Self:
        """
        Update an existing repository in the database.

        :param repository_service: the repository service to use
        :param id: the id of the repository
        :param name: the name of the repository
        :param description: the description of the repository
        :param metadata: the metadata of the request
        :return: the updated repository
        """
        response = repository_service.update_repository(
            id, name, description, metadata=metadata
        )
        return cls(
            id=response.id,
            name=response.name,
            description=response.description,
            last_updated=response.last_update,
        )

    @classmethod
    def get(
        cls,
        repository_service: RepositoryService,
        id: str,
        metadata: MetadataTuple,
    ) -> Self:
        """
        Get an existing repository from the database.

        :param repository_service: the repository service to use
        :param id: the id of the repository
        :param metadata: the metadata of the request
        :return: the repository
        """
        response = repository_service.get_repository(id, metadata=metadata)
        return cls(
            id=response.id,
            name=response.name,
            description=response.description,
            last_updated=response.last_update,
        )

    @classmethod
    def get_by_name(
        cls,
        settings: Settings,
        repository_service: RepositoryService,
        name: str,
        metadata: MetadataTuple,
    ) -> Self:
        """
        Get a repository by its name.

        :param settings: the settings to use
        :param repository_service: the repository service to use
        :param name: the name of the repository
        :param metadata: the metadata of the request
        :return: the repository with the given name
        """
        response = repository_service.get_by_name(name, metadata=metadata)
        obj = cls(
            id=response.id,
            name=response.name,
            description=response.description,
            last_updated=response.last_update,
        )
        obj.get_clone_url(settings)
        return obj

    @classmethod
    def delete(
        cls,
        repository_service: RepositoryService,
        id: str,
        metadata: MetadataTuple,
    ) -> None:
        """
        Delete an existing repository from the database.

        :param repository_service: the repository service to use
        :param id: the id of the repository
        :param metadata: the metadata of the request
        :return: None
        """
        repository_service.delete_repository(id, metadata=metadata)

    @classmethod
    def get_all(
        cls, repository_service: RepositoryService, metadata: MetadataTuple
    ) -> list[Self]:
        """
        Get all existing repositories from the database.

        :param repository_service: the repository service to use
        :param metadata: the metadata of the request
        :return: a list of all existing repositories
        """
        response_list = repository_service.get_repositories(metadata=metadata)
        return [
            cls(
                id=response.id,
                name=response.name,
                description=response.description,
                last_updated=response.last_update,
            )
            for response in response_list.repositories
        ]
