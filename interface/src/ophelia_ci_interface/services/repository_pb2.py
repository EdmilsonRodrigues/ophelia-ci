# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: repository.proto
# Protobuf Python Version: 5.29.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    0,
    '',
    'repository.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
import ophelia_ci_interface.services.common_pb2 as common__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x10repository.proto\x12\nrepository\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x0c\x63ommon.proto\"0\n\x14GetRepositoryRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0c\n\x04name\x18\x02 \x01(\t\"O\n\x17\x43reateRepositoryRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x02 \x01(\t\x12\x11\n\tgitignore\x18\x03 \x01(\t\"H\n\x17UpdateRepositoryRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\"%\n\x17\x44\x65leteRepositoryRequest\x12\n\n\x02id\x18\x01 \x01(\t\"t\n\x12RepositoryResponse\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\x12/\n\x0blast_update\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\"N\n\x16ListRepositoryResponse\x12\x34\n\x0crepositories\x18\x01 \x03(\x0b\x32\x1e.repository.RepositoryResponse2\xa5\x03\n\x11RepositoryService\x12W\n\x10\x43reateRepository\x12#.repository.CreateRepositoryRequest\x1a\x1e.repository.RepositoryResponse\x12W\n\x10UpdateRepository\x12#.repository.UpdateRepositoryRequest\x1a\x1e.repository.RepositoryResponse\x12\x43\n\x0eListRepository\x12\r.common.Empty\x1a\".repository.ListRepositoryResponse\x12Q\n\rGetRepository\x12 .repository.GetRepositoryRequest\x1a\x1e.repository.RepositoryResponse\x12\x46\n\x10\x44\x65leteRepository\x12#.repository.DeleteRepositoryRequest\x1a\r.common.EmptyB)Z\'github.com/EdmilsonRodrigues/ophelia-cib\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'repository_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\'github.com/EdmilsonRodrigues/ophelia-ci'
  _globals['_GETREPOSITORYREQUEST']._serialized_start=79
  _globals['_GETREPOSITORYREQUEST']._serialized_end=127
  _globals['_CREATEREPOSITORYREQUEST']._serialized_start=129
  _globals['_CREATEREPOSITORYREQUEST']._serialized_end=208
  _globals['_UPDATEREPOSITORYREQUEST']._serialized_start=210
  _globals['_UPDATEREPOSITORYREQUEST']._serialized_end=282
  _globals['_DELETEREPOSITORYREQUEST']._serialized_start=284
  _globals['_DELETEREPOSITORYREQUEST']._serialized_end=321
  _globals['_REPOSITORYRESPONSE']._serialized_start=323
  _globals['_REPOSITORYRESPONSE']._serialized_end=439
  _globals['_LISTREPOSITORYRESPONSE']._serialized_start=441
  _globals['_LISTREPOSITORYRESPONSE']._serialized_end=519
  _globals['_REPOSITORYSERVICE']._serialized_start=522
  _globals['_REPOSITORYSERVICE']._serialized_end=943
# @@protoc_insertion_point(module_scope)
