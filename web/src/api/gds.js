const api = require('./trisa/gds/api/v1beta1/api_grpc_web_pb');
const models = require('./trisa/gds/models/v1beta1/models_pb');

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
        req.setRegisteredDirectory("vaspdirectory.net");
        break
      case "common name":
        req.setCommonName(query);
        break
      default:
        throw new Error("unacceptable input type to lookup query");
    }

    let self = this;
    return new Promise((resolve, reject) => {
      self.client.lookup(req, {}, (err, rep) => {
        if (err || !rep) {
          reject(err);
        } else {
          resolve(rep.toObject());
        }
      });
    });
  }

  status = (commonName) => {
    const req = new api.StatusRequest();
    req.setCommonName(commonName);
    this.client.status(req, {}, (err, rep) => {
      if (err || !rep) {
        console.log(err);
        return
      }
      console.log(this.verificationStatus(rep.getVerificationStatus()));
      console.log(rep.toObject());
    });
  };

  verificationStatus = (num) => {
    for (const key in models.VerificationState) {
      if (num === models.VerificationState[key]) {
        return key;
      }
    };
  }

}


export default new GDS(defaultEndpoint());