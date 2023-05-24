import axiosInstance from 'utils/axios';
import { NetworkType } from 'types/enums';
export const registrationRequest = async (network: NetworkType, body: any) => {
  const response = await axiosInstance.post(`/register/${network}`, { ...body });
  return response.data;
};
