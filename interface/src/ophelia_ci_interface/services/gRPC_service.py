import base64
import hashlib

import grpc
import ophelia_ci_interface.services.repository_pb2 as repository_pb2
import ophelia_ci_interface.services.repository_pb2_grpc as repository_pb2_grpc
import ophelia_ci_interface.services.user_pb2 as user_pb2
import ophelia_ci_interface.services.user_pb2_grpc as user_pb2_grpc
import paramiko


class RepositoryService:
    def __init__(self, server: str):
        self.server = server
        self.channel = None
        self.stub = None

    def __enter__(self):
        self.channel = grpc.insecure_channel(self.server)
        self.stub = repository_pb2_grpc.RepositoryServiceStub(self.channel)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.channel.close()

    def get_status(self):
        return 'Connected'

    def get_repositories(self, metadata: tuple[tuple[str, str]]):
        with self:
            response_list = self.stub.ListRepository(
                repository_pb2.Empty(), metadata=metadata
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
                    metadata=metadata,
                )
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
                    metadata=metadata,
                )
            )
        return response_update

    def get_repository(self, id: str, metadata: tuple[tuple[str, str]]):
        with self:
            response_get = self.stub.GetRepository(
                repository_pb2.GetRepositoryRequest(id=id, metadata=metadata)
            )
        return response_get

    def get_by_name(self, name: str, metadata: tuple[tuple[str, str]]):
        with self:
            response_get = self.stub.GetRepository(
                repository_pb2.GetRepositoryRequest(
                    name=name, metadata=metadata
                )
            )
        return response_get

    def delete_repository(self, id: str, metadata: tuple[tuple[str, str]]):
        with self:
            response_delete = self.stub.DeleteRepository(
                repository_pb2.DeleteRepositoryRequest(
                    id=id, metadata=metadata
                )
            )
        return response_delete


class AuthenticationService:
    def __init__(self, server: str):
        self.server = server
        self.channel = None
        self.stub = None

    def __enter__(self):
        self.channel = grpc.insecure_channel(self.server)
        self.stub = user_pb2_grpc.AuthenticationStub(self.channel)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.channel.close()
        self.stub = None

    def authenticate(self, username: str, private_key: str):
        with self:
            challenge = self.request_challenge(username)
            private_key_obj = paramiko.RSAKey(file_obj=private_key)
            challenge_bytes = base64.b64decode(challenge.challenge)
            hash_obj = hashlib.sha256(challenge_bytes)
            signature = private_key_obj.sign_ssh_data(hash_obj.digest())
            signature_b64 = base64.b64encode(signature).decode('utf-8')
            response_auth = self.stub.Authentication(
                user_pb2.AuthenticationRequest(
                    username=username, signature=signature_b64
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


class UserService:
    def __init__(self, server: str):
        self.server = server
        self.channel = None
        self.stub = None

    def __enter__(self):
        self.channel = grpc.insecure_channel(self.server)
        self.stub = user_pb2_grpc.UserServiceStub(self.channel)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.channel.close()
        self.stub = None

    def create_user(
        self, username: str, private_key: str, metadata: tuple[tuple[str, str]]
    ):
        with self:
            response_create = self.stub.CreateUser(
                user_pb2.CreateUserRequest(
                    username=username,
                    private_key=private_key,
                    metadata=metadata,
                )
            )
        return response_create

    def get_user(self, id: str, metadata: tuple[tuple[str, str]]):
        with self:
            response_get = self.stub.GetUser(
                user_pb2.GetUserRequest(id=id, metadata=metadata)
            )
        return response_get

    def get_user_by_username(
        self, username: str, metadata: tuple[tuple[str, str]]
    ):
        with self:
            response_get = self.stub.GetUser(
                user_pb2.GetUserRequest(username=username, metadata=metadata)
            )
        return response_get

    def get_users(self, metadata: tuple[tuple[str, str]]):
        with self:
            response_get = self.stub.ListUser(
                user_pb2.GetUsersRequest(metadata=metadata)
            )
        return response_get

    def update_user(
        self,
        id: str,
        username: str,
        private_key: str,
        metadata: tuple[tuple[str, str]],
    ):
        with self:
            response_update = self.stub.UpdateUser(
                user_pb2.UpdateUserRequest(
                    id=id,
                    username=username,
                    private_key=private_key,
                    metadata=metadata,
                )
            )
        return response_update

    def delete_user(self, id: str, metadata: tuple[tuple[str, str]]):
        with self:
            response_delete = self.stub.DeleteUser(
                user_pb2.DeleteUserRequest(id=id, metadata=metadata)
            )
        return response_delete
