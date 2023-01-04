import { APICore } from '@/helpers/api/apiCore';

const api = new APICore();

const getAllAutocomplete = () => api.get('/autocomplete');

export default getAllAutocomplete;
