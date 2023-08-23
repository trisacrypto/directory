import axiosInstance from 'utils/axios';

export const lookup = async (query: string) => {
  const response = await axiosInstance.get(`/lookup?${query}`);
  return response.data;
};

export const lookupAutocomplete = async () => {
  const response = await axiosInstance.get(`/lookup/autocomplete`);
  return response.data;
};
