import axiosInstance from 'utils/axios';

export const getCertificateStep = async (network: string, key: any) => {
  const response = await axiosInstance.get(`/register/${network}?step=${key}`);
  return response.data;
};
