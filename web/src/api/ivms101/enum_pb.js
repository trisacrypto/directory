// source: ivms101/enum.proto
/**
 * @fileoverview
 * @enhanceable
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

goog.exportSymbol('proto.ivms101.AddressTypeCode', null, global);
goog.exportSymbol('proto.ivms101.LegalPersonNameTypeCode', null, global);
goog.exportSymbol('proto.ivms101.NationalIdentifierTypeCode', null, global);
goog.exportSymbol('proto.ivms101.NaturalPersonNameTypeCode', null, global);
goog.exportSymbol('proto.ivms101.TransliterationMethodCode', null, global);
/**
 * @enum {number}
 */
proto.ivms101.NaturalPersonNameTypeCode = {
  NATURAL_PERSON_NAME_TYPE_CODE_ALIA: 0,
  NATURAL_PERSON_NAME_TYPE_CODE_BIRT: 1,
  NATURAL_PERSON_NAME_TYPE_CODE_MAID: 2,
  NATURAL_PERSON_NAME_TYPE_CODE_LEGL: 3,
  NATURAL_PERSON_NAME_TYPE_CODE_MISC: 4
};

/**
 * @enum {number}
 */
proto.ivms101.LegalPersonNameTypeCode = {
  LEGAL_PERSON_NAME_TYPE_CODE_LEGL: 0,
  LEGAL_PERSON_NAME_TYPE_CODE_SHRT: 1,
  LEGAL_PERSON_NAME_TYPE_CODE_TRAD: 2
};

/**
 * @enum {number}
 */
proto.ivms101.AddressTypeCode = {
  ADDRESS_TYPE_CODE_HOME: 0,
  ADDRESS_TYPE_CODE_BIZZ: 1,
  ADDRESS_TYPE_CODE_GEOG: 2
};

/**
 * @enum {number}
 */
proto.ivms101.NationalIdentifierTypeCode = {
  NATIONAL_IDENTIFIER_TYPE_CODE_ARNU: 0,
  NATIONAL_IDENTIFIER_TYPE_CODE_CCPT: 1,
  NATIONAL_IDENTIFIER_TYPE_CODE_RAID: 2,
  NATIONAL_IDENTIFIER_TYPE_CODE_DRLC: 3,
  NATIONAL_IDENTIFIER_TYPE_CODE_FIIN: 4,
  NATIONAL_IDENTIFIER_TYPE_CODE_TXID: 5,
  NATIONAL_IDENTIFIER_TYPE_CODE_SOCS: 6,
  NATIONAL_IDENTIFIER_TYPE_CODE_IDCD: 7,
  NATIONAL_IDENTIFIER_TYPE_CODE_LEIX: 8,
  NATIONAL_IDENTIFIER_TYPE_CODE_MISC: 9
};

/**
 * @enum {number}
 */
proto.ivms101.TransliterationMethodCode = {
  TRANSLITERATION_METHOD_CODE_ARAB: 0,
  TRANSLITERATION_METHOD_CODE_ARAN: 1,
  TRANSLITERATION_METHOD_CODE_ARMN: 2,
  TRANSLITERATION_METHOD_CODE_CYRL: 3,
  TRANSLITERATION_METHOD_CODE_DEVA: 4,
  TRANSLITERATION_METHOD_CODE_GEOR: 5,
  TRANSLITERATION_METHOD_CODE_GREK: 6,
  TRANSLITERATION_METHOD_CODE_HANI: 7,
  TRANSLITERATION_METHOD_CODE_HEBR: 8,
  TRANSLITERATION_METHOD_CODE_KANA: 10,
  TRANSLITERATION_METHOD_CODE_KORE: 11,
  TRANSLITERATION_METHOD_CODE_THAI: 12,
  TRANSLITERATION_METHOD_CODE_OTHR: 13
};

goog.object.extend(exports, proto.ivms101);
