import { t } from "@lingui/macro";


const _ = require('lodash');
const api = require('./trisa/gds/api/v1beta1/api_grpc_web_pb');
const models = require('./trisa/gds/models/v1beta1/models_pb');
const ivms101 = require('./ivms101/ivms101_pb');

const registeredDirectory = "trisa.directory";

const defaultEndpoint = () => {
  // Use environment configured variable by default
  if (process.env.REACT_APP_GDS_API_ENDPOINT) {
    return process.env.REACT_APP_GDS_API_ENDPOINT;
  }

  // Make an educated guess based on the node environment
  switch (process.env.NODE_ENV) {
    case "development":
      return "http://localhost:8080"
    case "production":
      return "https://proxy.trisa.directory"
    default:
      throw new Error(t`could not identify default GDS api endpoint`);
  }
}

// HACK: type identification by field name for recusrive proto message construction.
const protoTypeMap = {
  trixo: models.TRIXOQuestionnaire,
  contacts: models.Contacts,
  technical: models.Contact,
  administrative: models.Contact,
  billing: models.Contact,
  legal: models.Contact,
  entity: ivms101.LegalPerson,
  name: ivms101.LegalPersonName,
  name_identifiers: ivms101.LegalPersonNameId,
  local_name_identifiers: ivms101.LocalLegalPersonNameId,
  phonetic_name_identifiers: ivms101.LocalLegalPersonNameId,
  geographic_addresses: ivms101.Address,
  national_identification: ivms101.NationalIdentification,
  other_jurisdictions: models.Jurisdiction
};

// This is kind of a nightmare ...
// See: https://github.com/grpc/grpc-web/issues/875
const protoFromObject = (ProtoClass, obj) => {
  let msg = new ProtoClass();
  for (const field in obj) {
    let value = obj[field];
    let setter = `set${_.upperFirst(_.camelCase(field))}`;

    if (!_.isArray(value) && !_.isObject(value)) {
      // Handle Primitives
      if (msg[setter]) {
        msg[setter](value);
      } else {
        throw new Error(t`field ${field} with setter ${setter} does not exist on ${ProtoClass}`);
      }
    } else if (_.isArray(value)) {
      // Handle Repeated
      setter = `add${_.upperFirst(_.camelCase(field))}`;
      if (msg[setter]) {
        for (const item of value) {
          if (_.isObject(item)) {
            // Handle nested objects in repeated arrays
            const NestedProto = protoTypeMap[field];
            if (NestedProto) {
              const nested = protoFromObject(NestedProto, item);
              msg[setter](nested);
            } else {
              throw new Error(t`unknown nested proto type for field ${field} with setter ${setter}`);
            }
          } else {
            // Add primitive inside repeated array
            msg[setter](item);
          }
        }
      } else {
        throw new Error(t`repeated field ${field} with setter ${setter} does not exist on ${ProtoClass}`);
      }
    } else if (_.isObject(value)) {
      // Recursively call protoFromObject
      // The issue is that we don't know the type, so using custom type map ...
      const NestedProto = protoTypeMap[field];
      if (NestedProto) {
        const nested = protoFromObject(NestedProto, value);
        if (msg[setter]) {
          msg[setter](nested);
        } else {
          throw new Error(t`field ${field} with setter ${setter} does not exist on ${ProtoClass}`);
        }
      } else {
        throw new Error(t`unknown nested proto type for field ${field} with setter ${setter}`);
      }
    } else {
      throw new Error(t`unhandled field ${field} with setter ${setter} and value type ${typeof(value)}`);
    }
  }
  return msg
}


class GDS {
  constructor(endpoint) {
    console.log("accessing GDS API at", endpoint)
    this.client = new api.TRISADirectoryClient(endpoint);
  }

  // Utilize the Lookup RPC to query for a single VASP record either by UUID or by
  // common name depending on the inputType.
  lookup = (query, inputType) => {
    const req = new api.LookupRequest();
    switch (inputType) {
      case "uuid":
        req.setId(query);
        req.setRegisteredDirectory(registeredDirectory);
        break
      case "common name":
        req.setCommonName(query);
        break
      default:
        throw new Error(t`unacceptable input type to lookup query`);
    }

    let client = this.client;
    return new Promise((resolve, reject) => {
      client.lookup(req, {}, (err, rep) => {
        if (err || !rep) {
          reject(err);
        } else {
          resolve(rep.toObject());
        }
      });
    });
  }

  search = (query) => {
    // let client = this.client;
    return new Promise((resolve, reject) => {
      reject(new Error(t`search API not implemented yet`));
    })
  }

  register = (formData) => {
    // Create protocol buffer
    const req = protoFromObject(api.RegisterRequest, formData);

    // Server requires contact to be null if it is not set.
    for (const contact in formData.contacts) {
      const value = formData.contacts[contact]
      if (!value.name && !value.email && !value.phone) {
        let contacts = req.getContacts();
        let setter = `set${_.upperFirst(_.camelCase(contact))}`;
        if (contacts[setter]) {
          contacts[setter](null)
          req.setContacts(contacts);
        } else {
          throw new Error(t`could not nullify empty ${contact} contact`);
        }
      }
    }

    console.log(req.toObject());
    let client = this.client;
    return new Promise((resolve, reject) => {
      client.register(req, {}, (err, rep) => {
        if (err || !rep) {
          reject(err);
        } else {
          let data = rep.toObject();
          data.status = this.verificationStatus(rep.getStatus());
          resolve(data);
        }
      });
    })
  }

  verifyContact = (vaspID, token) => {
    if (!vaspID || !token) {
      throw new Error(t`vaspID and token are required`);
    }

    const req = new api.VerifyContactRequest();
    req.setId(vaspID);
    req.setToken(token);

    let client = this.client;
    let verificationStatus = this.verificationStatus;
    return new Promise((resolve, reject) => {
      client.verifyContact(req, {}, (err, rep) => {
        if (err || !rep) {
          reject(err);
        } else {
          let data = rep.toObject();
          data.status = verificationStatus(rep.getStatus());
          resolve(data);
        }
      });
    });
  }

  verification = (commonName) => {
    const req = new api.VerificationRequest();
    req.setCommonName(commonName);

    let client = this.client;
    let verificationStatus = this.verificationStatus;
    return new Promise((resolve, reject) => {
      client.verification(req, {}, (err, rep) => {
        if (err || !rep) {
          reject(err);
        } else {
          console.log(verificationStatus(rep.getVerificationStatus()));
          resolve(rep.toObject());
        }
      });
    });
  };

  verificationStatus = (num) => {
    for (const key in models.VerificationState) {
      if (num === models.VerificationState[key]) {
        return key;
      }
    };
  }

  status = () => {
    const req = new api.HealthCheck();
    let client = this.client;
    return new Promise((resolve, reject) => {
      client.status(req, {}, (err, rep) => {
        if (err || !rep) {
          reject(err);
        } else {
          resolve(rep.toObject());
        }
      });
    });
  }

}


export default new GDS(defaultEndpoint());