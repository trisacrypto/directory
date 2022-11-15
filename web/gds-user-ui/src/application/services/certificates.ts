import axiosInstance from 'utils/axios';

export const getCertificates = async () => {
  const response = await axiosInstance.get('/certificates');
  return response.data;
};
