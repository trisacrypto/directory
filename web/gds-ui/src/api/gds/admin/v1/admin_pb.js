// source: gds/admin/v1/admin.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

var trisa_gds_api_v1beta1_api_pb = require('../../../trisa/gds/api/v1beta1/api_pb.js');
goog.object.extend(proto, trisa_gds_api_v1beta1_api_pb);
var trisa_gds_models_v1beta1_models_pb = require('../../../trisa/gds/models/v1beta1/models_pb.js');
goog.object.extend(proto, trisa_gds_models_v1beta1_models_pb);
goog.exportSymbol('proto.gds.admin.v1.ResendReply', null, global);
goog.exportSymbol('proto.gds.admin.v1.ResendRequest', null, global);
goog.exportSymbol('proto.gds.admin.v1.ResendRequest.ResendType', null, global);
goog.exportSymbol('proto.gds.admin.v1.ReviewReply', null, global);
goog.exportSymbol('proto.gds.admin.v1.ReviewRequest', null, global);
goog.exportSymbol('proto.gds.admin.v1.StatusReply', null, global);
goog.exportSymbol('proto.gds.admin.v1.StatusRequest', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.gds.admin.v1.ReviewRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gds.admin.v1.ReviewRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.gds.admin.v1.ReviewRequest.displayName = 'proto.gds.admin.v1.ReviewRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.gds.admin.v1.ReviewReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gds.admin.v1.ReviewReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.gds.admin.v1.ReviewReply.displayName = 'proto.gds.admin.v1.ReviewReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.gds.admin.v1.ResendRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gds.admin.v1.ResendRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.gds.admin.v1.ResendRequest.displayName = 'proto.gds.admin.v1.ResendRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.gds.admin.v1.ResendReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gds.admin.v1.ResendReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.gds.admin.v1.ResendReply.displayName = 'proto.gds.admin.v1.ResendReply';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.gds.admin.v1.StatusRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gds.admin.v1.StatusRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.gds.admin.v1.StatusRequest.displayName = 'proto.gds.admin.v1.StatusRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.gds.admin.v1.StatusReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gds.admin.v1.StatusReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.gds.admin.v1.StatusReply.displayName = 'proto.gds.admin.v1.StatusReply';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.gds.admin.v1.ReviewRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.gds.admin.v1.ReviewRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gds.admin.v1.ReviewRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.ReviewRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    adminVerificationToken: jspb.Message.getFieldWithDefault(msg, 2, ""),
    accept: jspb.Message.getBooleanFieldWithDefault(msg, 3, false),
    rejectReason: jspb.Message.getFieldWithDefault(msg, 4, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.gds.admin.v1.ReviewRequest}
 */
proto.gds.admin.v1.ReviewRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gds.admin.v1.ReviewRequest;
  return proto.gds.admin.v1.ReviewRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gds.admin.v1.ReviewRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gds.admin.v1.ReviewRequest}
 */
proto.gds.admin.v1.ReviewRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAdminVerificationToken(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setAccept(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setRejectReason(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.gds.admin.v1.ReviewRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gds.admin.v1.ReviewRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gds.admin.v1.ReviewRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.ReviewRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAdminVerificationToken();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getAccept();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
  f = message.getRejectReason();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.gds.admin.v1.ReviewRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.gds.admin.v1.ReviewRequest} returns this
 */
proto.gds.admin.v1.ReviewRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string admin_verification_token = 2;
 * @return {string}
 */
proto.gds.admin.v1.ReviewRequest.prototype.getAdminVerificationToken = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.gds.admin.v1.ReviewRequest} returns this
 */
proto.gds.admin.v1.ReviewRequest.prototype.setAdminVerificationToken = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bool accept = 3;
 * @return {boolean}
 */
proto.gds.admin.v1.ReviewRequest.prototype.getAccept = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.gds.admin.v1.ReviewRequest} returns this
 */
proto.gds.admin.v1.ReviewRequest.prototype.setAccept = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
};


/**
 * optional string reject_reason = 4;
 * @return {string}
 */
proto.gds.admin.v1.ReviewRequest.prototype.getRejectReason = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.gds.admin.v1.ReviewRequest} returns this
 */
proto.gds.admin.v1.ReviewRequest.prototype.setRejectReason = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.gds.admin.v1.ReviewReply.prototype.toObject = function(opt_includeInstance) {
  return proto.gds.admin.v1.ReviewReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gds.admin.v1.ReviewReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.ReviewReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    error: (f = msg.getError()) && trisa_gds_api_v1beta1_api_pb.Error.toObject(includeInstance, f),
    status: jspb.Message.getFieldWithDefault(msg, 2, 0),
    message: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.gds.admin.v1.ReviewReply}
 */
proto.gds.admin.v1.ReviewReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gds.admin.v1.ReviewReply;
  return proto.gds.admin.v1.ReviewReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gds.admin.v1.ReviewReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gds.admin.v1.ReviewReply}
 */
proto.gds.admin.v1.ReviewReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new trisa_gds_api_v1beta1_api_pb.Error;
      reader.readMessage(value,trisa_gds_api_v1beta1_api_pb.Error.deserializeBinaryFromReader);
      msg.setError(value);
      break;
    case 2:
      var value = /** @type {!proto.trisa.gds.models.v1beta1.VerificationState} */ (reader.readEnum());
      msg.setStatus(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setMessage(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.gds.admin.v1.ReviewReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gds.admin.v1.ReviewReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gds.admin.v1.ReviewReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.ReviewReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getError();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      trisa_gds_api_v1beta1_api_pb.Error.serializeBinaryToWriter
    );
  }
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getMessage();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional trisa.gds.api.v1beta1.Error error = 1;
 * @return {?proto.trisa.gds.api.v1beta1.Error}
 */
proto.gds.admin.v1.ReviewReply.prototype.getError = function() {
  return /** @type{?proto.trisa.gds.api.v1beta1.Error} */ (
    jspb.Message.getWrapperField(this, trisa_gds_api_v1beta1_api_pb.Error, 1));
};


/**
 * @param {?proto.trisa.gds.api.v1beta1.Error|undefined} value
 * @return {!proto.gds.admin.v1.ReviewReply} returns this
*/
proto.gds.admin.v1.ReviewReply.prototype.setError = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.gds.admin.v1.ReviewReply} returns this
 */
proto.gds.admin.v1.ReviewReply.prototype.clearError = function() {
  return this.setError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.gds.admin.v1.ReviewReply.prototype.hasError = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional trisa.gds.models.v1beta1.VerificationState status = 2;
 * @return {!proto.trisa.gds.models.v1beta1.VerificationState}
 */
proto.gds.admin.v1.ReviewReply.prototype.getStatus = function() {
  return /** @type {!proto.trisa.gds.models.v1beta1.VerificationState} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.trisa.gds.models.v1beta1.VerificationState} value
 * @return {!proto.gds.admin.v1.ReviewReply} returns this
 */
proto.gds.admin.v1.ReviewReply.prototype.setStatus = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional string message = 3;
 * @return {string}
 */
proto.gds.admin.v1.ReviewReply.prototype.getMessage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.gds.admin.v1.ReviewReply} returns this
 */
proto.gds.admin.v1.ReviewReply.prototype.setMessage = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.gds.admin.v1.ResendRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.gds.admin.v1.ResendRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gds.admin.v1.ResendRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.ResendRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    type: jspb.Message.getFieldWithDefault(msg, 2, 0),
    reason: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.gds.admin.v1.ResendRequest}
 */
proto.gds.admin.v1.ResendRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gds.admin.v1.ResendRequest;
  return proto.gds.admin.v1.ResendRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gds.admin.v1.ResendRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gds.admin.v1.ResendRequest}
 */
proto.gds.admin.v1.ResendRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {!proto.gds.admin.v1.ResendRequest.ResendType} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setReason(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.gds.admin.v1.ResendRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gds.admin.v1.ResendRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gds.admin.v1.ResendRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.ResendRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getReason();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.gds.admin.v1.ResendRequest.ResendType = {
  UNKNOWN: 0,
  VERIFY_CONTACT: 1,
  REVIEW: 2,
  DELIVER_CERTS: 3,
  REJECTION: 4
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.gds.admin.v1.ResendRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.gds.admin.v1.ResendRequest} returns this
 */
proto.gds.admin.v1.ResendRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ResendType type = 2;
 * @return {!proto.gds.admin.v1.ResendRequest.ResendType}
 */
proto.gds.admin.v1.ResendRequest.prototype.getType = function() {
  return /** @type {!proto.gds.admin.v1.ResendRequest.ResendType} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.gds.admin.v1.ResendRequest.ResendType} value
 * @return {!proto.gds.admin.v1.ResendRequest} returns this
 */
proto.gds.admin.v1.ResendRequest.prototype.setType = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional string reason = 3;
 * @return {string}
 */
proto.gds.admin.v1.ResendRequest.prototype.getReason = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.gds.admin.v1.ResendRequest} returns this
 */
proto.gds.admin.v1.ResendRequest.prototype.setReason = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.gds.admin.v1.ResendReply.prototype.toObject = function(opt_includeInstance) {
  return proto.gds.admin.v1.ResendReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gds.admin.v1.ResendReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.ResendReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    sent: jspb.Message.getFieldWithDefault(msg, 1, 0),
    message: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.gds.admin.v1.ResendReply}
 */
proto.gds.admin.v1.ResendReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gds.admin.v1.ResendReply;
  return proto.gds.admin.v1.ResendReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gds.admin.v1.ResendReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gds.admin.v1.ResendReply}
 */
proto.gds.admin.v1.ResendReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setSent(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setMessage(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.gds.admin.v1.ResendReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gds.admin.v1.ResendReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gds.admin.v1.ResendReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.ResendReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSent();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getMessage();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional int64 sent = 1;
 * @return {number}
 */
proto.gds.admin.v1.ResendReply.prototype.getSent = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.gds.admin.v1.ResendReply} returns this
 */
proto.gds.admin.v1.ResendReply.prototype.setSent = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string message = 2;
 * @return {string}
 */
proto.gds.admin.v1.ResendReply.prototype.getMessage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.gds.admin.v1.ResendReply} returns this
 */
proto.gds.admin.v1.ResendReply.prototype.setMessage = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.gds.admin.v1.StatusRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.gds.admin.v1.StatusRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gds.admin.v1.StatusRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.StatusRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    noRegistrations: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    noCertificateRequests: jspb.Message.getBooleanFieldWithDefault(msg, 2, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.gds.admin.v1.StatusRequest}
 */
proto.gds.admin.v1.StatusRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gds.admin.v1.StatusRequest;
  return proto.gds.admin.v1.StatusRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gds.admin.v1.StatusRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gds.admin.v1.StatusRequest}
 */
proto.gds.admin.v1.StatusRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setNoRegistrations(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setNoCertificateRequests(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.gds.admin.v1.StatusRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gds.admin.v1.StatusRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gds.admin.v1.StatusRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.StatusRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getNoRegistrations();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getNoCertificateRequests();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
};


/**
 * optional bool no_registrations = 1;
 * @return {boolean}
 */
proto.gds.admin.v1.StatusRequest.prototype.getNoRegistrations = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.gds.admin.v1.StatusRequest} returns this
 */
proto.gds.admin.v1.StatusRequest.prototype.setNoRegistrations = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional bool no_certificate_requests = 2;
 * @return {boolean}
 */
proto.gds.admin.v1.StatusRequest.prototype.getNoCertificateRequests = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 2, false));
};


/**
 * @param {boolean} value
 * @return {!proto.gds.admin.v1.StatusRequest} returns this
 */
proto.gds.admin.v1.StatusRequest.prototype.setNoCertificateRequests = function(value) {
  return jspb.Message.setProto3BooleanField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.gds.admin.v1.StatusReply.prototype.toObject = function(opt_includeInstance) {
  return proto.gds.admin.v1.StatusReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gds.admin.v1.StatusReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.StatusReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    registrationsMap: (f = msg.getRegistrationsMap()) ? f.toObject(includeInstance, undefined) : [],
    certificateRequestsMap: (f = msg.getCertificateRequestsMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.gds.admin.v1.StatusReply}
 */
proto.gds.admin.v1.StatusReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gds.admin.v1.StatusReply;
  return proto.gds.admin.v1.StatusReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gds.admin.v1.StatusReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gds.admin.v1.StatusReply}
 */
proto.gds.admin.v1.StatusReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = msg.getRegistrationsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readInt64, null, "", 0);
         });
      break;
    case 2:
      var value = msg.getCertificateRequestsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readInt64, null, "", 0);
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.gds.admin.v1.StatusReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gds.admin.v1.StatusReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gds.admin.v1.StatusReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gds.admin.v1.StatusReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRegistrationsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(1, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeInt64);
  }
  f = message.getCertificateRequestsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(2, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeInt64);
  }
};


/**
 * map<string, int64> registrations = 1;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,number>}
 */
proto.gds.admin.v1.StatusReply.prototype.getRegistrationsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,number>} */ (
      jspb.Message.getMapField(this, 1, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.gds.admin.v1.StatusReply} returns this
 */
proto.gds.admin.v1.StatusReply.prototype.clearRegistrationsMap = function() {
  this.getRegistrationsMap().clear();
  return this;};


/**
 * map<string, int64> certificate_requests = 2;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,number>}
 */
proto.gds.admin.v1.StatusReply.prototype.getCertificateRequestsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,number>} */ (
      jspb.Message.getMapField(this, 2, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.gds.admin.v1.StatusReply} returns this
 */
proto.gds.admin.v1.StatusReply.prototype.clearCertificateRequestsMap = function() {
  this.getCertificateRequestsMap().clear();
  return this;};


goog.object.extend(exports, proto.gds.admin.v1);
