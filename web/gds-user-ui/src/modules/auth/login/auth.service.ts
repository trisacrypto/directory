import axiosInstance from 'utils/axios';
import { getCookie } from 'utils/cookies';
export const logUserInBff = async () => {
  const response = await axiosInstance.post(
    `/users/login`,
    {},
    {
      headers: {
        Authorization: `Bearer ${getCookie('access_token')}`
      }
    }
  );

  return response;
};
export const getUserRoles = async () => {
  const response = await axiosInstance.get(`/users/roles`, {
    headers: {
      Authorization: `Bearer ${getCookie('access_token')}`
    }
  });

  return response;
};
