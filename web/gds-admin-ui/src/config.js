const config = {
    ENVIRONMENT: process.env.NODE_ENV,
    GOOGLE_CLIENT_ID: process.env.REACT_APP_GOOGLE_CLIENT_ID,
    IS_TESTNET: JSON.parse(process.env.REACT_APP_GDS_IS_TESTNET),
    GDS_API_URL: process.env.REACT_APP_GDS_API_ENDPOINT,
    REACT_APP_VERSION_NUMBER: process.env.REACT_APP_VERSION_NUMBER,
    REACT_APP_GIT_REVISION: process.env.REACT_APP_GIT_REVISION,
    get isProd() {
        return process.env.NODE_ENV === 'production'
    },
    get appVersion() {
        const self = process.env

        if (self.REACT_APP_VERSION_NUMBER) {
            return self.REACT_APP_VERSION_NUMBER
        }

        return undefined
    },
    get gitVersion() {
        const self = process.env

        if (self.REACT_APP_GIT_REVISION) {
            return self.REACT_APP_GIT_REVISION
        }

        return undefined
    }
};

export default config;
