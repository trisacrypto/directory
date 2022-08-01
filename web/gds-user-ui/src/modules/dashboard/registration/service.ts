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

// submit testnet registration endpoint request

export const submitTestnetRegistration = async () => {
  setAuthorization();
  const response = await axiosInstance.post(`/register/testnet`);
  return response;
};

// submit mainnet registration endpoint request

export const submitMainnetRegistration = async () => {
  setAuthorization();
  const response = await axiosInstance.post(`/register/mainnet`);
  return response;
};
