import { APICore } from '@/helpers/api/apiCore';

const api = new APICore();

export default function getAppVersion() {
  return api.get('/status');
}
