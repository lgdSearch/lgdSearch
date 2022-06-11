from pymilvus import connections
from pymilvus import CollectionSchema, FieldSchema, DataType
from pymilvus import utility
from pymilvus import Collection
import os

def connect2Milvus():
    connections.connect(
        alias="default",
        host='192.168.242.5',
        port='19530'
    )


def create_collection():
    connect2Milvus()
    if not utility.has_collection("image"):
        img_id = FieldSchema(
            name="img_id",
            dtype=DataType.INT64,
            is_primary=True,
        )
        featrure = FieldSchema(
            name="feature",
            dtype=DataType.FLOAT_VECTOR,
            dim=2048
        )
        schema = CollectionSchema(
            fields=[img_id, featrure],
            description="Test img search"
        )
        collection_name = "image"
        Collection(
            name=collection_name,
            schema=schema,
            using='default',
            shards_num=2,
            consistency_level="Strong"
        )


def get_collection_image():
    connect2Milvus()
    if not utility.has_collection("image"):
        img_id = FieldSchema(
            name="img_id",
            dtype=DataType.INT64,
            is_primary=True,
        )
        featrure = FieldSchema(
            name="feature",
            dtype=DataType.FLOAT_VECTOR,
            dim=2048
        )
        schema = CollectionSchema(
            fields=[img_id, featrure],
            description="Test img search"
        )
        collection_name = "image"
        collection = Collection(
            name=collection_name,
            schema=schema,
            using='default',
            consistency_level="Strong"
        )
        index_params = {
            "metric_type": "L2",
            "index_type": "IVF_FLAT",
            "params": {"nlist": 16384}
        }
        collection.create_index(
            field_name="feature",
            index_params=index_params
        )
        return collection
    return Collection("image")


def insert(data):
    collection = get_collection_image()
    collection.insert(data)


def search(feat):
    collection = get_collection_image()
    collection.load()
    search_params = {"metric_type": "L2", "params": {"nprobe": 16}}
    results = collection.search(
        data=[feat],
        anns_field="feature",
        param=search_params,
        limit=10,
        consistency_level="Strong"
    )
    collection.release()
    return results[0].ids


def drop_image():
    connect2Milvus()
    if utility.has_collection("image"):
        utility.drop_collection("image")


if __name__ == '__main__':
    drop_image()