# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import get_feature_pb2 as get__feature__pb2


class GrpcServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.getFeature = channel.unary_unary(
                '/GrpcService/getFeature',
                request_serializer=get__feature__pb2.Request.SerializeToString,
                response_deserializer=get__feature__pb2.Response.FromString,
                )


class GrpcServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def getFeature(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_GrpcServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'getFeature': grpc.unary_unary_rpc_method_handler(
                    servicer.getFeature,
                    request_deserializer=get__feature__pb2.Request.FromString,
                    response_serializer=get__feature__pb2.Response.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'GrpcService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class GrpcService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def getFeature(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/GrpcService/getFeature',
            get__feature__pb2.Request.SerializeToString,
            get__feature__pb2.Response.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
