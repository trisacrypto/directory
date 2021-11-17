/**
 * @fileoverview gRPC-Web generated client stub for gds.admin.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var trisa_gds_api_v1beta1_api_pb = require('../../../trisa/gds/api/v1beta1/api_pb.js')

var trisa_gds_models_v1beta1_models_pb = require('../../../trisa/gds/models/v1beta1/models_pb.js')
const proto = {};
proto.gds = {};
proto.gds.admin = {};
proto.gds.admin.v1 = require('./admin_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.gds.admin.v1.DirectoryAdministrationClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.gds.admin.v1.DirectoryAdministrationPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gds.admin.v1.ReviewRequest,
 *   !proto.gds.admin.v1.ReviewReply>}
 */
const methodDescriptor_DirectoryAdministration_Review = new grpc.web.MethodDescriptor(
  '/gds.admin.v1.DirectoryAdministration/Review',
  grpc.web.MethodType.UNARY,
  proto.gds.admin.v1.ReviewRequest,
  proto.gds.admin.v1.ReviewReply,
  /**
   * @param {!proto.gds.admin.v1.ReviewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.gds.admin.v1.ReviewReply.deserializeBinary
);


/**
 * @param {!proto.gds.admin.v1.ReviewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gds.admin.v1.ReviewReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gds.admin.v1.ReviewReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gds.admin.v1.DirectoryAdministrationClient.prototype.review =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gds.admin.v1.DirectoryAdministration/Review',
      request,
      metadata || {},
      methodDescriptor_DirectoryAdministration_Review,
      callback);
};


/**
 * @param {!proto.gds.admin.v1.ReviewRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gds.admin.v1.ReviewReply>}
 *     Promise that resolves to the response
 */
proto.gds.admin.v1.DirectoryAdministrationPromiseClient.prototype.review =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gds.admin.v1.DirectoryAdministration/Review',
      request,
      metadata || {},
      methodDescriptor_DirectoryAdministration_Review);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gds.admin.v1.ResendRequest,
 *   !proto.gds.admin.v1.ResendReply>}
 */
const methodDescriptor_DirectoryAdministration_Resend = new grpc.web.MethodDescriptor(
  '/gds.admin.v1.DirectoryAdministration/Resend',
  grpc.web.MethodType.UNARY,
  proto.gds.admin.v1.ResendRequest,
  proto.gds.admin.v1.ResendReply,
  /**
   * @param {!proto.gds.admin.v1.ResendRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.gds.admin.v1.ResendReply.deserializeBinary
);


/**
 * @param {!proto.gds.admin.v1.ResendRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gds.admin.v1.ResendReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gds.admin.v1.ResendReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gds.admin.v1.DirectoryAdministrationClient.prototype.resend =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gds.admin.v1.DirectoryAdministration/Resend',
      request,
      metadata || {},
      methodDescriptor_DirectoryAdministration_Resend,
      callback);
};


/**
 * @param {!proto.gds.admin.v1.ResendRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gds.admin.v1.ResendReply>}
 *     Promise that resolves to the response
 */
proto.gds.admin.v1.DirectoryAdministrationPromiseClient.prototype.resend =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gds.admin.v1.DirectoryAdministration/Resend',
      request,
      metadata || {},
      methodDescriptor_DirectoryAdministration_Resend);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gds.admin.v1.StatusRequest,
 *   !proto.gds.admin.v1.StatusReply>}
 */
const methodDescriptor_DirectoryAdministration_Status = new grpc.web.MethodDescriptor(
  '/gds.admin.v1.DirectoryAdministration/Status',
  grpc.web.MethodType.UNARY,
  proto.gds.admin.v1.StatusRequest,
  proto.gds.admin.v1.StatusReply,
  /**
   * @param {!proto.gds.admin.v1.StatusRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.gds.admin.v1.StatusReply.deserializeBinary
);


/**
 * @param {!proto.gds.admin.v1.StatusRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gds.admin.v1.StatusReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gds.admin.v1.StatusReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gds.admin.v1.DirectoryAdministrationClient.prototype.status =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gds.admin.v1.DirectoryAdministration/Status',
      request,
      metadata || {},
      methodDescriptor_DirectoryAdministration_Status,
      callback);
};


/**
 * @param {!proto.gds.admin.v1.StatusRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gds.admin.v1.StatusReply>}
 *     Promise that resolves to the response
 */
proto.gds.admin.v1.DirectoryAdministrationPromiseClient.prototype.status =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gds.admin.v1.DirectoryAdministration/Status',
      request,
      metadata || {},
      methodDescriptor_DirectoryAdministration_Status);
};


module.exports = proto.gds.admin.v1;

