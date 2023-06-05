// @generated by protoc-gen-es v1.2.1 with parameter "target=ts"
// @generated from file rill/runtime/v1/connectors.proto (package rill.runtime.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, protoInt64, Timestamp } from "@bufbuild/protobuf";

/**
 * @generated from message rill.runtime.v1.S3Object
 */
export class S3Object extends Message<S3Object> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: google.protobuf.Timestamp modified_on = 2;
   */
  modifiedOn?: Timestamp;

  /**
   * @generated from field: int64 size = 3;
   */
  size = protoInt64.zero;

  /**
   * @generated from field: bool is_dir = 4;
   */
  isDir = false;

  constructor(data?: PartialMessage<S3Object>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.S3Object";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "modified_on", kind: "message", T: Timestamp },
    { no: 3, name: "size", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
    { no: 4, name: "is_dir", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): S3Object {
    return new S3Object().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): S3Object {
    return new S3Object().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): S3Object {
    return new S3Object().fromJsonString(jsonString, options);
  }

  static equals(a: S3Object | PlainMessage<S3Object> | undefined, b: S3Object | PlainMessage<S3Object> | undefined): boolean {
    return proto3.util.equals(S3Object, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.S3ListBucketsRequest
 */
export class S3ListBucketsRequest extends Message<S3ListBucketsRequest> {
  /**
   * @generated from field: uint32 page_size = 1;
   */
  pageSize = 0;

  /**
   * @generated from field: string page_token = 2;
   */
  pageToken = "";

  constructor(data?: PartialMessage<S3ListBucketsRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.S3ListBucketsRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "page_size", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "page_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): S3ListBucketsRequest {
    return new S3ListBucketsRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): S3ListBucketsRequest {
    return new S3ListBucketsRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): S3ListBucketsRequest {
    return new S3ListBucketsRequest().fromJsonString(jsonString, options);
  }

  static equals(a: S3ListBucketsRequest | PlainMessage<S3ListBucketsRequest> | undefined, b: S3ListBucketsRequest | PlainMessage<S3ListBucketsRequest> | undefined): boolean {
    return proto3.util.equals(S3ListBucketsRequest, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.S3ListBucketsResponse
 */
export class S3ListBucketsResponse extends Message<S3ListBucketsResponse> {
  /**
   * @generated from field: string next_page_token = 1;
   */
  nextPageToken = "";

  /**
   * @generated from field: repeated string buckets = 2;
   */
  buckets: string[] = [];

  constructor(data?: PartialMessage<S3ListBucketsResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.S3ListBucketsResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "next_page_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "buckets", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): S3ListBucketsResponse {
    return new S3ListBucketsResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): S3ListBucketsResponse {
    return new S3ListBucketsResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): S3ListBucketsResponse {
    return new S3ListBucketsResponse().fromJsonString(jsonString, options);
  }

  static equals(a: S3ListBucketsResponse | PlainMessage<S3ListBucketsResponse> | undefined, b: S3ListBucketsResponse | PlainMessage<S3ListBucketsResponse> | undefined): boolean {
    return proto3.util.equals(S3ListBucketsResponse, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.S3ListObjectsRequest
 */
export class S3ListObjectsRequest extends Message<S3ListObjectsRequest> {
  /**
   * @generated from field: uint32 page_size = 1;
   */
  pageSize = 0;

  /**
   * @generated from field: string page_token = 2;
   */
  pageToken = "";

  /**
   * @generated from field: string bucket = 3;
   */
  bucket = "";

  /**
   * @generated from field: string region = 4;
   */
  region = "";

  /**
   * @generated from field: string prefix = 5;
   */
  prefix = "";

  /**
   * @generated from field: string start_after = 6;
   */
  startAfter = "";

  /**
   * @generated from field: string delimiter = 7;
   */
  delimiter = "";

  constructor(data?: PartialMessage<S3ListObjectsRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.S3ListObjectsRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "page_size", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "page_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "bucket", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "region", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "prefix", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "start_after", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 7, name: "delimiter", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): S3ListObjectsRequest {
    return new S3ListObjectsRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): S3ListObjectsRequest {
    return new S3ListObjectsRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): S3ListObjectsRequest {
    return new S3ListObjectsRequest().fromJsonString(jsonString, options);
  }

  static equals(a: S3ListObjectsRequest | PlainMessage<S3ListObjectsRequest> | undefined, b: S3ListObjectsRequest | PlainMessage<S3ListObjectsRequest> | undefined): boolean {
    return proto3.util.equals(S3ListObjectsRequest, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.S3ListObjectsResponse
 */
export class S3ListObjectsResponse extends Message<S3ListObjectsResponse> {
  /**
   * @generated from field: string next_page_token = 1;
   */
  nextPageToken = "";

  /**
   * @generated from field: repeated rill.runtime.v1.S3Object objects = 2;
   */
  objects: S3Object[] = [];

  constructor(data?: PartialMessage<S3ListObjectsResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.S3ListObjectsResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "next_page_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "objects", kind: "message", T: S3Object, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): S3ListObjectsResponse {
    return new S3ListObjectsResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): S3ListObjectsResponse {
    return new S3ListObjectsResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): S3ListObjectsResponse {
    return new S3ListObjectsResponse().fromJsonString(jsonString, options);
  }

  static equals(a: S3ListObjectsResponse | PlainMessage<S3ListObjectsResponse> | undefined, b: S3ListObjectsResponse | PlainMessage<S3ListObjectsResponse> | undefined): boolean {
    return proto3.util.equals(S3ListObjectsResponse, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.S3GetBucketMetadataRequest
 */
export class S3GetBucketMetadataRequest extends Message<S3GetBucketMetadataRequest> {
  /**
   * @generated from field: string bucket = 1;
   */
  bucket = "";

  constructor(data?: PartialMessage<S3GetBucketMetadataRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.S3GetBucketMetadataRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "bucket", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): S3GetBucketMetadataRequest {
    return new S3GetBucketMetadataRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): S3GetBucketMetadataRequest {
    return new S3GetBucketMetadataRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): S3GetBucketMetadataRequest {
    return new S3GetBucketMetadataRequest().fromJsonString(jsonString, options);
  }

  static equals(a: S3GetBucketMetadataRequest | PlainMessage<S3GetBucketMetadataRequest> | undefined, b: S3GetBucketMetadataRequest | PlainMessage<S3GetBucketMetadataRequest> | undefined): boolean {
    return proto3.util.equals(S3GetBucketMetadataRequest, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.S3GetBucketMetadataResponse
 */
export class S3GetBucketMetadataResponse extends Message<S3GetBucketMetadataResponse> {
  /**
   * @generated from field: string region = 1;
   */
  region = "";

  constructor(data?: PartialMessage<S3GetBucketMetadataResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.S3GetBucketMetadataResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "region", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): S3GetBucketMetadataResponse {
    return new S3GetBucketMetadataResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): S3GetBucketMetadataResponse {
    return new S3GetBucketMetadataResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): S3GetBucketMetadataResponse {
    return new S3GetBucketMetadataResponse().fromJsonString(jsonString, options);
  }

  static equals(a: S3GetBucketMetadataResponse | PlainMessage<S3GetBucketMetadataResponse> | undefined, b: S3GetBucketMetadataResponse | PlainMessage<S3GetBucketMetadataResponse> | undefined): boolean {
    return proto3.util.equals(S3GetBucketMetadataResponse, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.S3GetCredentialsInfoRequest
 */
export class S3GetCredentialsInfoRequest extends Message<S3GetCredentialsInfoRequest> {
  constructor(data?: PartialMessage<S3GetCredentialsInfoRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.S3GetCredentialsInfoRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): S3GetCredentialsInfoRequest {
    return new S3GetCredentialsInfoRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): S3GetCredentialsInfoRequest {
    return new S3GetCredentialsInfoRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): S3GetCredentialsInfoRequest {
    return new S3GetCredentialsInfoRequest().fromJsonString(jsonString, options);
  }

  static equals(a: S3GetCredentialsInfoRequest | PlainMessage<S3GetCredentialsInfoRequest> | undefined, b: S3GetCredentialsInfoRequest | PlainMessage<S3GetCredentialsInfoRequest> | undefined): boolean {
    return proto3.util.equals(S3GetCredentialsInfoRequest, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.S3GetCredentialsInfoResponse
 */
export class S3GetCredentialsInfoResponse extends Message<S3GetCredentialsInfoResponse> {
  /**
   * @generated from field: bool exist = 1;
   */
  exist = false;

  /**
   * @generated from field: string provider = 2;
   */
  provider = "";

  constructor(data?: PartialMessage<S3GetCredentialsInfoResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.S3GetCredentialsInfoResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "exist", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 2, name: "provider", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): S3GetCredentialsInfoResponse {
    return new S3GetCredentialsInfoResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): S3GetCredentialsInfoResponse {
    return new S3GetCredentialsInfoResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): S3GetCredentialsInfoResponse {
    return new S3GetCredentialsInfoResponse().fromJsonString(jsonString, options);
  }

  static equals(a: S3GetCredentialsInfoResponse | PlainMessage<S3GetCredentialsInfoResponse> | undefined, b: S3GetCredentialsInfoResponse | PlainMessage<S3GetCredentialsInfoResponse> | undefined): boolean {
    return proto3.util.equals(S3GetCredentialsInfoResponse, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.GCSObject
 */
export class GCSObject extends Message<GCSObject> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: google.protobuf.Timestamp modified_on = 2;
   */
  modifiedOn?: Timestamp;

  /**
   * @generated from field: int64 size = 3;
   */
  size = protoInt64.zero;

  /**
   * @generated from field: bool is_dir = 4;
   */
  isDir = false;

  constructor(data?: PartialMessage<GCSObject>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.GCSObject";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "modified_on", kind: "message", T: Timestamp },
    { no: 3, name: "size", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
    { no: 4, name: "is_dir", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GCSObject {
    return new GCSObject().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GCSObject {
    return new GCSObject().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GCSObject {
    return new GCSObject().fromJsonString(jsonString, options);
  }

  static equals(a: GCSObject | PlainMessage<GCSObject> | undefined, b: GCSObject | PlainMessage<GCSObject> | undefined): boolean {
    return proto3.util.equals(GCSObject, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.GCSListBucketsRequest
 */
export class GCSListBucketsRequest extends Message<GCSListBucketsRequest> {
  /**
   * @generated from field: uint32 page_size = 1;
   */
  pageSize = 0;

  /**
   * @generated from field: string page_token = 2;
   */
  pageToken = "";

  constructor(data?: PartialMessage<GCSListBucketsRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.GCSListBucketsRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "page_size", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "page_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GCSListBucketsRequest {
    return new GCSListBucketsRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GCSListBucketsRequest {
    return new GCSListBucketsRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GCSListBucketsRequest {
    return new GCSListBucketsRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GCSListBucketsRequest | PlainMessage<GCSListBucketsRequest> | undefined, b: GCSListBucketsRequest | PlainMessage<GCSListBucketsRequest> | undefined): boolean {
    return proto3.util.equals(GCSListBucketsRequest, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.GCSListBucketsResponse
 */
export class GCSListBucketsResponse extends Message<GCSListBucketsResponse> {
  /**
   * @generated from field: string next_page_token = 1;
   */
  nextPageToken = "";

  /**
   * @generated from field: repeated string buckets = 2;
   */
  buckets: string[] = [];

  constructor(data?: PartialMessage<GCSListBucketsResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.GCSListBucketsResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "next_page_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "buckets", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GCSListBucketsResponse {
    return new GCSListBucketsResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GCSListBucketsResponse {
    return new GCSListBucketsResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GCSListBucketsResponse {
    return new GCSListBucketsResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GCSListBucketsResponse | PlainMessage<GCSListBucketsResponse> | undefined, b: GCSListBucketsResponse | PlainMessage<GCSListBucketsResponse> | undefined): boolean {
    return proto3.util.equals(GCSListBucketsResponse, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.GCSListObjectsRequest
 */
export class GCSListObjectsRequest extends Message<GCSListObjectsRequest> {
  /**
   * @generated from field: uint32 page_size = 1;
   */
  pageSize = 0;

  /**
   * @generated from field: string page_token = 2;
   */
  pageToken = "";

  /**
   * @generated from field: string bucket = 3;
   */
  bucket = "";

  /**
   * @generated from field: string prefix = 4;
   */
  prefix = "";

  /**
   * @generated from field: string start_offset = 5;
   */
  startOffset = "";

  /**
   * @generated from field: string end_offset = 6;
   */
  endOffset = "";

  /**
   * @generated from field: string delimiter = 7;
   */
  delimiter = "";

  constructor(data?: PartialMessage<GCSListObjectsRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.GCSListObjectsRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "page_size", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "page_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "bucket", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "prefix", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "start_offset", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "end_offset", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 7, name: "delimiter", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GCSListObjectsRequest {
    return new GCSListObjectsRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GCSListObjectsRequest {
    return new GCSListObjectsRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GCSListObjectsRequest {
    return new GCSListObjectsRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GCSListObjectsRequest | PlainMessage<GCSListObjectsRequest> | undefined, b: GCSListObjectsRequest | PlainMessage<GCSListObjectsRequest> | undefined): boolean {
    return proto3.util.equals(GCSListObjectsRequest, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.GCSListObjectsResponse
 */
export class GCSListObjectsResponse extends Message<GCSListObjectsResponse> {
  /**
   * @generated from field: string next_page_token = 1;
   */
  nextPageToken = "";

  /**
   * @generated from field: repeated rill.runtime.v1.GCSObject objects = 2;
   */
  objects: GCSObject[] = [];

  constructor(data?: PartialMessage<GCSListObjectsResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.GCSListObjectsResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "next_page_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "objects", kind: "message", T: GCSObject, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GCSListObjectsResponse {
    return new GCSListObjectsResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GCSListObjectsResponse {
    return new GCSListObjectsResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GCSListObjectsResponse {
    return new GCSListObjectsResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GCSListObjectsResponse | PlainMessage<GCSListObjectsResponse> | undefined, b: GCSListObjectsResponse | PlainMessage<GCSListObjectsResponse> | undefined): boolean {
    return proto3.util.equals(GCSListObjectsResponse, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.GCSGetCredentialsInfoRequest
 */
export class GCSGetCredentialsInfoRequest extends Message<GCSGetCredentialsInfoRequest> {
  constructor(data?: PartialMessage<GCSGetCredentialsInfoRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.GCSGetCredentialsInfoRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GCSGetCredentialsInfoRequest {
    return new GCSGetCredentialsInfoRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GCSGetCredentialsInfoRequest {
    return new GCSGetCredentialsInfoRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GCSGetCredentialsInfoRequest {
    return new GCSGetCredentialsInfoRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GCSGetCredentialsInfoRequest | PlainMessage<GCSGetCredentialsInfoRequest> | undefined, b: GCSGetCredentialsInfoRequest | PlainMessage<GCSGetCredentialsInfoRequest> | undefined): boolean {
    return proto3.util.equals(GCSGetCredentialsInfoRequest, a, b);
  }
}

/**
 * @generated from message rill.runtime.v1.GCSGetCredentialsInfoResponse
 */
export class GCSGetCredentialsInfoResponse extends Message<GCSGetCredentialsInfoResponse> {
  /**
   * @generated from field: bool exist = 1;
   */
  exist = false;

  /**
   * @generated from field: string project_id = 2;
   */
  projectId = "";

  constructor(data?: PartialMessage<GCSGetCredentialsInfoResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.GCSGetCredentialsInfoResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "exist", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 2, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GCSGetCredentialsInfoResponse {
    return new GCSGetCredentialsInfoResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GCSGetCredentialsInfoResponse {
    return new GCSGetCredentialsInfoResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GCSGetCredentialsInfoResponse {
    return new GCSGetCredentialsInfoResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GCSGetCredentialsInfoResponse | PlainMessage<GCSGetCredentialsInfoResponse> | undefined, b: GCSGetCredentialsInfoResponse | PlainMessage<GCSGetCredentialsInfoResponse> | undefined): boolean {
    return proto3.util.equals(GCSGetCredentialsInfoResponse, a, b);
  }
}
