import base64
import hashlib
import io
import logging

import grpc  # type: ignore[import-untyped]  # type: ignore[import-untyped]  # type: ignore[import-untyped]  # type: ignore[import-untyped]
import ophelia_ci_interface.services.common_pb2 as common_pb2
import ophelia_ci_interface.services.health_pb2_grpc as health_pb2_grpc
import ophelia_ci_interface.services.repository_pb2 as repository_pb2
import ophelia_ci_interface.services.repository_pb2_grpc as repository_pb2_grpc
import ophelia_ci_interface.services.user_pb2 as user_pb2
import ophelia_ci_interface.services.user_pb2_grpc as user_pb2_grpc
import paramiko
from ophelia_ci_interface.models.generals import (
    OpheliaException,
    log_formatted,
)


class ServiceMixin[T]:
    """
    Mixin class for gRPC services. Implements the context manager protocol.
    The context manager is used to open and close the gRPC channel.


    :param server: the server address
    :param stub_class: the stub class
    :param stub: the stub instance
    """

    server: str
    stub_class: type[T]
    stub: T

    def __init__(self, server: str):
        self.server = server
        self.channel = None

    def open_channel(self):
        self.channel = grpc.insecure_channel(self.server)

    def __enter__(self):
        self.open_channel()
        self.stub = self.stub_class(self.channel)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.channel.close()


class HealthService(ServiceMixin):
    """
    Health service for checking the status of the server.

    :param server: the server address
    :param stub_class: the stub class
    :param stub: the stub instance
    """

    stub_class = health_pb2_grpc.HealthServiceStub

    def get_status(self) -> str:
        """
        Check the status of the server by invoking the Health RPC.

        :return: 'Connected' if the server responds successfully,
                otherwise 'Failed Connecting' if an exception occurs.
        """
        with self:
            try:
                self.stub.Health(common_pb2.Empty())
                return 'Connected'
            except grpc.RpcError:
                log_formatted('Failed Connecting', logging_level=logging.ERROR)
                return 'Failed Connecting'


class RepositoryService(ServiceMixin):
    stub_class = repository_pb2_grpc.RepositoryServiceStub

    def get_repositories(self, metadata: tuple[tuple[str, str]]):
        """
        Retrieve a list of all repositories from the database.

        :param metadata: The metadata for the request.
        :return: A response containing a list of repositories.
        """
        with self:
            log_formatted('Getting repositories')
            response_list = self.stub.ListRepository(
                common_pb2.Empty(), metadata=metadata
            )
        return response_list

    def create_repository(
        self,
        name: str,
        description: str,
        gitignore: str,
        metadata: tuple[tuple[str, str]],
    ):
        """
        Create a new repository in the database.

        :param name: the name of the repository
        :param description: the description of the repository
        :param gitignore: the main language of the repository, to be used for
            the base gitignore
        :param metadata: the metadata of the request
        :return: the newly created repository
        """
        with self:
            log_formatted(
                'Creating repository',
                name=name,
                description=description,
                gitignore=gitignore,
            )
            response_create = self.stub.CreateRepository(
                repository_pb2.CreateRepositoryRequest(
                    name=name,
                    description=description,
                    gitignore=gitignore,
                ),
                metadata=metadata,
            )
        return response_create

    def update_repository(
        self,
        id: str,
        name: str,
        description: str,
        metadata: tuple[tuple[str, str]],
    ):
        """
        Update an existing repository in the database.

        :param id: the ID of the repository
        :param name: the name of the repository
        :param description: the description of the repository
        :param metadata: the metadata of the request
        :return: the updated repository
        """
        with self:
            log_formatted(
                'Updating repository',
                id=id,
                name=name,
                description=description,
            )
            response_update = self.stub.UpdateRepository(
                repository_pb2.UpdateRepositoryRequest(
                    id=id,
                    name=name,
                    description=description,
                ),
                metadata=metadata,
            )
        return response_update

    def get_repository(self, id: str, metadata: tuple[tuple[str, str]]):
        """
        Retrieve a repository by its ID from the database.

        :param id: the ID of the repository
        :param metadata: the metadata for the request
        :return: the retrieved repository
        """
        with self:
            log_formatted('Getting repository', id=id)
            response_get = self.stub.GetRepository(
                repository_pb2.GetRepositoryRequest(id=id), metadata=metadata
            )
        return response_get

    def get_by_name(self, name: str, metadata: tuple[tuple[str, str]]):
        """
        Retrieve a repository by its name from the database.

        :param name: The name of the repository to retrieve.
        :param metadata: The metadata of the request.
        :return: The repository response from the database.
        """
        with self:
            log_formatted('Getting repository', name=name)
            response_get = self.stub.GetRepository(
                repository_pb2.GetRepositoryRequest(name=name),
                metadata=metadata,
            )
        return response_get

    def delete_repository(self, id: str, metadata: tuple[tuple[str, str]]):
        """
        Delete an existing repository from the database.

        :param id: the id of the repository
        :param metadata: the metadata of the request
        :return: the deleted repository
        """
        with self:
            log_formatted('Deleting repository', id=id)
            response_delete = self.stub.DeleteRepository(
                repository_pb2.DeleteRepositoryRequest(id=id),
                metadata=metadata,
            )
        return response_delete


class AuthenticationService(ServiceMixin):
    stub_class = user_pb2_grpc.AuthServiceStub

    def authenticate(self, username: str, private_key: str) -> str:
        """
        Authenticate a user with their username and private key.

        :param username: the username of the user
        :param private_key: the private key of the user
        :return: the token of the user if the authentication is successful
        :raises OpheliaException: if the authentication failed
        """
        with self:
            log_formatted('Authenticating user', username=username)
            challenge = self._request_challenge(username)
            private_key_obj = paramiko.RSAKey(
                file_obj=io.StringIO(private_key)
            )
            challenge_bytes = base64.b64decode(challenge.challenge)
            h = hashlib.sha256(challenge_bytes).digest()

            signature = private_key_obj.sign_ssh_data(h)

            blob = self.get_signature_blob(signature)

            # Send just the blob part
            signature_b64 = base64.b64encode(blob).decode('utf-8')

            log_formatted(
                'Signing challenge',
                challenge=challenge.challenge,
                signature_b64=signature_b64,
                logging_level=logging.DEBUG,
                sensitive=True,
            )

            response_auth = self.stub.Authentication(
                user_pb2.AuthenticationRequest(
                    username=username, challenge=signature_b64
                )
            )
            if response_auth.authenticated:
                return response_auth.token
            raise OpheliaException('Authentication failed')

    def _request_challenge(self, username: str):
        """
        Requests a challenge from the server to be signed by the user.

        :param username: the username of the user
        :return: the challenge to be signed
        """
        response_challenge = self.stub.AuthenticationChallenge(
            user_pb2.AuthenticationChallengeRequest(username=username)
        )
        return response_challenge

    @staticmethod
    def get_signature_blob(signature: paramiko.Message) -> bytes:
        """
        Get the signature blob from the SSH signature.

        SSH signature format:
        - 4 bytes: length of format string
        - format string (e.g., "ssh-rsa")
        - 4 bytes: length of signature blob
        - signature blob

        :param signature: the SSH signature
        :return: the signature blob

        """
        # Skip format string and length
        sig_io = io.BytesIO(bytes(signature))
        fmt_len = int.from_bytes(sig_io.read(4), byteorder='big')
        sig_io.read(fmt_len)  # Skip the format string

        # Read blob length and blob
        blob_len = int.from_bytes(sig_io.read(4), byteorder='big')
        return sig_io.read(blob_len)

    def authenticate_with_unique_key(self, unique_key: str):
        """
        Authenticate a user with the unique key created by the
        server on startup.

        :param unique_key: the unique key created by the server on startup
        :return: the token of the user if the authentication is successful
        :raises OpheliaException: if the authentication failed
        """
        with self:
            log_formatted(
                'Authenticating via Unique Key',
                unique_key=unique_key,
                logging_level=logging.WARNING,
            )
            response_auth = self.stub.UniqueKeyLogin(
                user_pb2.UniqueKeyLoginRequest(uniqueKey=unique_key)
            )
            if response_auth.authenticated:
                return response_auth.token
            raise OpheliaException('Authentication failed')


class UserService(ServiceMixin):
    stub_class = user_pb2_grpc.UserServiceStub

    def create_user(
        self, username: str, public_key: str, metadata: tuple[tuple[str, str]]
    ):
        """
        Create a new user in the database.

        :param username: the username of the user
        :param public_key: the public key of the user
        :param metadata: the metadata of the request
        :return: the newly created user
        """
        with self:
            log_formatted('Creating user', username=username)
            response_create = self.stub.CreateUser(
                user_pb2.CreateUserRequest(
                    username=username,
                    publicKey=public_key,
                ),
                metadata=metadata,
            )
        return response_create

    def get_user(self, id: str, metadata: tuple[tuple[str, str]]):
        """
        Get a user from the database by its id.

        :param id: the id of the user
        :param metadata: the metadata of the request
        :return: the user
        """
        with self:
            log_formatted('Getting user', id=id)
            response_get = self.stub.GetUser(
                user_pb2.GetUserRequest(id=id), metadata=metadata
            )
        return response_get

    def get_user_by_username(
        self, username: str, metadata: tuple[tuple[str, str]]
    ):
        """
        Get a user from the database by its username.

        :param username: the username of the user
        :param metadata: the metadata of the request
        :return: the user
        """
        with self:
            log_formatted('Getting user', username=username)
            response_get = self.stub.GetUser(
                user_pb2.GetUserRequest(username=username), metadata=metadata
            )
        return response_get

    def get_users(self, metadata: tuple[tuple[str, str]]):
        """
        Retrieve all users from the database.

        :param metadata: The metadata for the request.
        :return: A response containing a list of users.
        """
        with self:
            log_formatted('Getting users')
            response_get = self.stub.ListUser(
                common_pb2.Empty(), metadata=metadata
            )
        return response_get

    def update_user(
        self,
        id: str,
        username: str,
        public_key: str,
        metadata: tuple[tuple[str, str]],
    ):
        """
        Update an existing user in the database.

        :param id: the id of the user
        :param username: the username of the user
        :param public_key: the public key of the user
        :param metadata: the metadata of the request
        :return: the updated user
        """
        with self:
            log_formatted('Updating user', id=id, username=username)
            response_update = self.stub.UpdateUser(
                user_pb2.UpdateUserRequest(
                    id=id,
                    username=username,
                    publicKey=public_key,
                ),
                metadata=metadata,
            )
        return response_update

    def delete_user(self, id: str, metadata: tuple[tuple[str, str]]):
        """
        Delete an existing user from the database.

        :param id: the id of the user
        :param metadata: the metadata of the request
        :return: the deleted user
        """
        with self:
            log_formatted('Deleting user', id=id)
            response_delete = self.stub.DeleteUser(
                user_pb2.DeleteUserRequest(id=id), metadata=metadata
            )
        return response_delete
