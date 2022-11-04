import axiosInstance from 'utils/axios';
import { getCookie } from 'utils/cookies';
export const getAllCollaborators = async () => {
  const response = await axiosInstance.get(`/collaborators`, {
    headers: {
      Authorization: `Bearer ${getCookie('access_token')}`
    }
  });
  return response;
};
export const addCollaborator = async () => {
  const response = await axiosInstance.get(`/collaborators`, {
    headers: {
      Authorization: `Bearer ${getCookie('access_token')}`
    }
  });
  return response;
};
