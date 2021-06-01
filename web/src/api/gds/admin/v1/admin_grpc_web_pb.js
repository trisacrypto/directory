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
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.gds.admin.v1.DirectoryAdministrationClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

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
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.gds.admin.v1.DirectoryAdministrationPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

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
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.gds.admin.v1.ReviewRequest,
 *   !proto.gds.admin.v1.ReviewReply>}
 */
const methodInfo_DirectoryAdministration_Review = new grpc.web.AbstractClientBase.MethodInfo(
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
 * @param {function(?grpc.web.Error, ?proto.gds.admin.v1.ReviewReply)}
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
 * @param {?Object<string, string>} metadata User defined
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


module.exports = proto.gds.admin.v1;

