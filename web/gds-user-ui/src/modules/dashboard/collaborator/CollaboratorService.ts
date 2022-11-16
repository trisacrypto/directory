
import axiosInstance from 'utils/axios';
// import type { Collaborator } from 'components/Collaborators/CollaboratorType';

// import { getCookie } from 'utils/cookies';
export const getAllCollaborators = async () => {
  const response = await axiosInstance.get(`/collaborators`);
  return response;
};
export const createCollaborator = async (data: any): Promise<any> => {
  const response: any = await axiosInstance.post(`/collaborators`, {
    ...data
  });
  console.log('response', response);
  return response;
};

export const updateCollaborator = async (data: any): Promise<any> => {
  const response: any = await axiosInstance.put(`/collaborators`, {
    ...data
  });
  return response;
};

export const deleteCollaborator = async (id: string): Promise<any> => {
  const response: any = await axiosInstance.delete(`/collaborators/${id}`);
  return response;
};


