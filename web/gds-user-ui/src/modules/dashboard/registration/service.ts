import axiosInstance, { setAuthorization } from 'utils/axios';
import { getCookie } from 'utils/cookies';
export const getRegistrationData = async () => {
  setAuthorization();
  const response = await axiosInstance.get(`/register`);
  return response;
};
export const postRegistrationData = async (data: any) => {
  setAuthorization();
  const response = await axiosInstance.post(`/register`, { ...data });
  return response;
};
