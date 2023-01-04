import { getCookie } from '@/utils';

import { APICore } from './apiCore';

const api = new APICore();

function postCredentials(credentials) {
  const csrfToken = getCookie('csrf_token');
  return api.create('/authenticate', credentials, {
    headers: {
      'X-CSRF-TOKEN': csrfToken,
    },
  });
}

export { postCredentials };
