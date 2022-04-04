import axiosInstance from 'utils/axios';

export const registrationRequest = async (network: string, body: any) => {
  const response = await axiosInstance.post(`/register/${network}`, { ...body });
  return response.data;
};
