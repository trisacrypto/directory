import axiosInstance from 'utils/axios';

export const lookup = async (query: string) => {
  const response = await axiosInstance.get(`/lookup?${query}`);
  return response.data;
};
