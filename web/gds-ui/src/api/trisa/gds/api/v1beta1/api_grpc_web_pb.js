/**
 * @fileoverview gRPC-Web generated client stub for trisa.gds.api.v1beta1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var ivms101_ivms101_pb = require('../../../../ivms101/ivms101_pb.js')

var trisa_gds_models_v1beta1_models_pb = require('../../../../trisa/gds/models/v1beta1/models_pb.js')

var trisa_gds_models_v1beta1_ca_pb = require('../../../../trisa/gds/models/v1beta1/ca_pb.js')
const proto = {};
proto.trisa = {};
proto.trisa.gds = {};
proto.trisa.gds.api = {};
proto.trisa.gds.api.v1beta1 = require('./api_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryClient =
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
proto.trisa.gds.api.v1beta1.TRISADirectoryPromiseClient =
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
 *   !proto.trisa.gds.api.v1beta1.LookupRequest,
 *   !proto.trisa.gds.api.v1beta1.LookupReply>}
 */
const methodDescriptor_TRISADirectory_Lookup = new grpc.web.MethodDescriptor(
  '/trisa.gds.api.v1beta1.TRISADirectory/Lookup',
  grpc.web.MethodType.UNARY,
  proto.trisa.gds.api.v1beta1.LookupRequest,
  proto.trisa.gds.api.v1beta1.LookupReply,
  /**
   * @param {!proto.trisa.gds.api.v1beta1.LookupRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.trisa.gds.api.v1beta1.LookupReply.deserializeBinary
);


/**
 * @param {!proto.trisa.gds.api.v1beta1.LookupRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.trisa.gds.api.v1beta1.LookupReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.trisa.gds.api.v1beta1.LookupReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryClient.prototype.lookup =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/Lookup',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_Lookup,
      callback);
};


/**
 * @param {!proto.trisa.gds.api.v1beta1.LookupRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.trisa.gds.api.v1beta1.LookupReply>}
 *     Promise that resolves to the response
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryPromiseClient.prototype.lookup =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/Lookup',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_Lookup);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.trisa.gds.api.v1beta1.SearchRequest,
 *   !proto.trisa.gds.api.v1beta1.SearchReply>}
 */
const methodDescriptor_TRISADirectory_Search = new grpc.web.MethodDescriptor(
  '/trisa.gds.api.v1beta1.TRISADirectory/Search',
  grpc.web.MethodType.UNARY,
  proto.trisa.gds.api.v1beta1.SearchRequest,
  proto.trisa.gds.api.v1beta1.SearchReply,
  /**
   * @param {!proto.trisa.gds.api.v1beta1.SearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.trisa.gds.api.v1beta1.SearchReply.deserializeBinary
);


/**
 * @param {!proto.trisa.gds.api.v1beta1.SearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.trisa.gds.api.v1beta1.SearchReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.trisa.gds.api.v1beta1.SearchReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryClient.prototype.search =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/Search',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_Search,
      callback);
};


/**
 * @param {!proto.trisa.gds.api.v1beta1.SearchRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.trisa.gds.api.v1beta1.SearchReply>}
 *     Promise that resolves to the response
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryPromiseClient.prototype.search =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/Search',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_Search);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.trisa.gds.api.v1beta1.RegisterRequest,
 *   !proto.trisa.gds.api.v1beta1.RegisterReply>}
 */
const methodDescriptor_TRISADirectory_Register = new grpc.web.MethodDescriptor(
  '/trisa.gds.api.v1beta1.TRISADirectory/Register',
  grpc.web.MethodType.UNARY,
  proto.trisa.gds.api.v1beta1.RegisterRequest,
  proto.trisa.gds.api.v1beta1.RegisterReply,
  /**
   * @param {!proto.trisa.gds.api.v1beta1.RegisterRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.trisa.gds.api.v1beta1.RegisterReply.deserializeBinary
);


/**
 * @param {!proto.trisa.gds.api.v1beta1.RegisterRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.trisa.gds.api.v1beta1.RegisterReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.trisa.gds.api.v1beta1.RegisterReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryClient.prototype.register =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/Register',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_Register,
      callback);
};


/**
 * @param {!proto.trisa.gds.api.v1beta1.RegisterRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.trisa.gds.api.v1beta1.RegisterReply>}
 *     Promise that resolves to the response
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryPromiseClient.prototype.register =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/Register',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_Register);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.trisa.gds.api.v1beta1.VerifyContactRequest,
 *   !proto.trisa.gds.api.v1beta1.VerifyContactReply>}
 */
const methodDescriptor_TRISADirectory_VerifyContact = new grpc.web.MethodDescriptor(
  '/trisa.gds.api.v1beta1.TRISADirectory/VerifyContact',
  grpc.web.MethodType.UNARY,
  proto.trisa.gds.api.v1beta1.VerifyContactRequest,
  proto.trisa.gds.api.v1beta1.VerifyContactReply,
  /**
   * @param {!proto.trisa.gds.api.v1beta1.VerifyContactRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.trisa.gds.api.v1beta1.VerifyContactReply.deserializeBinary
);


/**
 * @param {!proto.trisa.gds.api.v1beta1.VerifyContactRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.trisa.gds.api.v1beta1.VerifyContactReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.trisa.gds.api.v1beta1.VerifyContactReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryClient.prototype.verifyContact =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/VerifyContact',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_VerifyContact,
      callback);
};


/**
 * @param {!proto.trisa.gds.api.v1beta1.VerifyContactRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.trisa.gds.api.v1beta1.VerifyContactReply>}
 *     Promise that resolves to the response
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryPromiseClient.prototype.verifyContact =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/VerifyContact',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_VerifyContact);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.trisa.gds.api.v1beta1.VerificationRequest,
 *   !proto.trisa.gds.api.v1beta1.VerificationReply>}
 */
const methodDescriptor_TRISADirectory_Verification = new grpc.web.MethodDescriptor(
  '/trisa.gds.api.v1beta1.TRISADirectory/Verification',
  grpc.web.MethodType.UNARY,
  proto.trisa.gds.api.v1beta1.VerificationRequest,
  proto.trisa.gds.api.v1beta1.VerificationReply,
  /**
   * @param {!proto.trisa.gds.api.v1beta1.VerificationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.trisa.gds.api.v1beta1.VerificationReply.deserializeBinary
);


/**
 * @param {!proto.trisa.gds.api.v1beta1.VerificationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.trisa.gds.api.v1beta1.VerificationReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.trisa.gds.api.v1beta1.VerificationReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryClient.prototype.verification =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/Verification',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_Verification,
      callback);
};


/**
 * @param {!proto.trisa.gds.api.v1beta1.VerificationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.trisa.gds.api.v1beta1.VerificationReply>}
 *     Promise that resolves to the response
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryPromiseClient.prototype.verification =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/Verification',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_Verification);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.trisa.gds.api.v1beta1.HealthCheck,
 *   !proto.trisa.gds.api.v1beta1.ServiceState>}
 */
const methodDescriptor_TRISADirectory_Status = new grpc.web.MethodDescriptor(
  '/trisa.gds.api.v1beta1.TRISADirectory/Status',
  grpc.web.MethodType.UNARY,
  proto.trisa.gds.api.v1beta1.HealthCheck,
  proto.trisa.gds.api.v1beta1.ServiceState,
  /**
   * @param {!proto.trisa.gds.api.v1beta1.HealthCheck} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.trisa.gds.api.v1beta1.ServiceState.deserializeBinary
);


/**
 * @param {!proto.trisa.gds.api.v1beta1.HealthCheck} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.trisa.gds.api.v1beta1.ServiceState)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.trisa.gds.api.v1beta1.ServiceState>|undefined}
 *     The XHR Node Readable Stream
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryClient.prototype.status =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/Status',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_Status,
      callback);
};


/**
 * @param {!proto.trisa.gds.api.v1beta1.HealthCheck} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.trisa.gds.api.v1beta1.ServiceState>}
 *     Promise that resolves to the response
 */
proto.trisa.gds.api.v1beta1.TRISADirectoryPromiseClient.prototype.status =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/trisa.gds.api.v1beta1.TRISADirectory/Status',
      request,
      metadata || {},
      methodDescriptor_TRISADirectory_Status);
};


module.exports = proto.trisa.gds.api.v1beta1;

