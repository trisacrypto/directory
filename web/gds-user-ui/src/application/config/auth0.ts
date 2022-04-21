import auth0 from 'auth0-js';
const defaultConfig: auth0.AuthOptions = {
  domain: process.env.REACT_APP_AUTH0_DOMAIN?.toString() || '',
  clientID: process.env.REACT_APP_AUTH0_CLIENT_ID?.toString() || '',
  redirectUri: process.env.REACT_APP_AUTH0_REDIRECT_URI,
  audience: process.env.REACT_APP_AUTH0_AUDIENCE,
  scope: process.env.REACT_APP_AUTH0_SCOPE || 'openid profile email',
  responseType: 'token id_token code'

  // responseType: 'code',
  // responseMode: 'query',
  // state: '',
  // nonce: '',
  // pkce: false,
  // display: 'page',
  // prompt: 'none',
  // maxAge: 0,
  // uiLocales: '',
  // claimsLocales: '',
  // idTokenHint: '',
  // loginHint: '',
  // acrValues: '',
  // resource: '',
};

const getAuth0Config = () => {
  return { ...defaultConfig };
};

export default getAuth0Config;
