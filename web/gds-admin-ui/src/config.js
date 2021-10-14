const config = {
    ENVIRONMENT: process.env.NODE_ENV,
    GOOGLE_CLIENT_ID: process.env.REACT_APP_GOOGLE_CLIENT_ID,
    IS_TESNET: JSON.parse(process.env.REACT_APP_GDS_IS_TESTNET),
    GDS_API_URL: process.env.REACT_APP_GDS_API_ENDPOINT
};

export default config;
