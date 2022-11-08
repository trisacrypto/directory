import axiosInstance from 'utils/axios';
// import { getCookie } from 'utils/cookies';
export const getAllCollaborators = async () => {
  const response = await axiosInstance.get(`/collaborators`);
  return response;
};
export const addCollaborator = async (data: any) => {
  const response = await axiosInstance.post(`/collaborators`, {
    ...data
  });
  return response;
};
