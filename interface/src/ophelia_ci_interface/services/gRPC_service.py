import base64
import hashlib
import io

import grpc  # type: ignore[import-untyped]  # type: ignore[import-untyped]  # type: ignore[import-untyped]  # type: ignore[import-untyped]
import ophelia_ci_interface.services.common_pb2 as common_pb2
import ophelia_ci_interface.services.health_pb2_grpc as health_pb2_grpc
import ophelia_ci_interface.services.repository_pb2 as repository_pb2
import ophelia_ci_interface.services.repository_pb2_grpc as repository_pb2_grpc
import ophelia_ci_interface.services.user_pb2 as user_pb2
import ophelia_ci_interface.services.user_pb2_grpc as user_pb2_grpc
import paramiko


class ServiceMixin[T]:
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
    stub_class = health_pb2_grpc.HealthServiceStub

    def get_status(self) -> str:
        with self:
            try:
                self.stub.Health(common_pb2.Empty())
                return 'Connected'
            except grpc.RpcError:
                return 'Failed Connecting'


class RepositoryService(ServiceMixin):
    stub_class = repository_pb2_grpc.RepositoryServiceStub

    def get_repositories(self, metadata: tuple[tuple[str, str]]):
        with self:
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
        with self:
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
        with self:
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
        with self:
            response_get = self.stub.GetRepository(
                repository_pb2.GetRepositoryRequest(id=id), metadata=metadata
            )
        return response_get

    def get_by_name(self, name: str, metadata: tuple[tuple[str, str]]):
        with self:
            response_get = self.stub.GetRepository(
                repository_pb2.GetRepositoryRequest(name=name),
                metadata=metadata,
            )
        return response_get

    def delete_repository(self, id: str, metadata: tuple[tuple[str, str]]):
        with self:
            response_delete = self.stub.DeleteRepository(
                repository_pb2.DeleteRepositoryRequest(id=id),
                metadata=metadata,
            )
        return response_delete


class AuthenticationService(ServiceMixin):
    stub_class = user_pb2_grpc.AuthServiceStub

    def authenticate(self, username: str, private_key: str):
        with self:
            challenge = self.request_challenge(username)
            private_key_obj = paramiko.RSAKey(
                file_obj=io.StringIO(private_key)
            )
            challenge_bytes = base64.b64decode(challenge.challenge)
            hash_obj = hashlib.sha256(challenge_bytes)
            signature = private_key_obj.sign_ssh_data(hash_obj.digest())
            signature_b64 = base64.b64encode(bytes(signature)).decode('utf-8')
            response_auth = self.stub.Authentication(
                user_pb2.AuthenticationRequest(
                    username=username, challenge=signature_b64
                )
            )
            if response_auth.authenticated:
                return response_auth.token
            raise Exception('Authentication failed')

    def request_challenge(self, username: str):
        response_challenge = self.stub.AuthenticationChallenge(
            user_pb2.AuthenticationChallengeRequest(username=username)
        )
        return response_challenge


class UserService(ServiceMixin):
    stub_class = user_pb2_grpc.UserServiceStub

    def create_user(
        self, username: str, public_key: str, metadata: tuple[tuple[str, str]]
    ):
        with self:
            response_create = self.stub.CreateUser(
                user_pb2.CreateUserRequest(
                    username=username,
                    publicKey=public_key,
                ),
                metadata=metadata,
            )
        return response_create

    def get_user(self, id: str, metadata: tuple[tuple[str, str]]):
        with self:
            response_get = self.stub.GetUser(
                user_pb2.GetUserRequest(id=id), metadata=metadata
            )
        return response_get

    def get_user_by_username(
        self, username: str, metadata: tuple[tuple[str, str]]
    ):
        with self:
            response_get = self.stub.GetUser(
                user_pb2.GetUserRequest(username=username), metadata=metadata
            )
        return response_get

    def get_users(self, metadata: tuple[tuple[str, str]]):
        with self:
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
        with self:
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
        with self:
            response_delete = self.stub.DeleteUser(
                user_pb2.DeleteUserRequest(id=id), metadata=metadata
            )
        return response_delete
