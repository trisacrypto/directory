// @flow
import { APICore } from './apiCore';

const api = new APICore();

function postCredentials(credentials, params) {
    return api.create('/authenticate', credentials, params)
}

export { postCredentials };
