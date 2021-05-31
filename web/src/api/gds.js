const api = require('./trisa/gds/api/v1beta1/api_grpc_web_pb');
const models = require('./trisa/gds/models/v1beta1/models_pb');

const registeredDirectory = "vaspdirectory.net";

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
      return "https://proxy.vaspdirectory.net"
    default:
      throw new Error("could not identify default GDS api endpoint");
  }
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
        throw new Error("unacceptable input type to lookup query");
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

  verifyContact = (vaspID, token) => {
    if (!vaspID || !token) {
      throw new Error("vaspID and token are required");
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