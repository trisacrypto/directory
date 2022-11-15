
import axiosInstance from 'utils/axios';
import type { Collaborator } from 'components/Collaborators/CollaboratorType';

// import { getCookie } from 'utils/cookies';
export const getAllCollaborators = async () => {
  const response = await axiosInstance.get(`/collaborators`);
  return response;
};
export const createCollaborator = async (data: any): Promise<Collaborator> => {
  const response: any = await axiosInstance.post(`/collaborators`, {
    ...data
  });
  console.log('response', response);
  return response;
};


