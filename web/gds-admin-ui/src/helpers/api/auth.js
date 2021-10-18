// @flow
import { getCookie } from '../../utils';
import { APICore } from './apiCore';

const api = new APICore();

function postCredentials(credentials, params) {
    return api.create('/authenticate', credentials, params)
}

function reauthenticate(credential, params) {
    const csrfToken = getCookie('csrf_token');

    return api.create('/reauthenticate', credential, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        },
        ...params
    })
}

export { postCredentials, reauthenticate };
