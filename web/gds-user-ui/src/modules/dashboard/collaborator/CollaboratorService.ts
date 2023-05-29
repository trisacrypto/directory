import axiosInstance from 'utils/axios';
// import type { Collaborator } from 'components/Collaborators/CollaboratorType';
interface TUpdateCollaborator {
  id: string;
  data: any;
}
// import { getCookie } from 'utils/cookies';
export const getAllCollaborators = async () => {
  const response = await axiosInstance.get(`/collaborators`);
  return response;
};
export const createCollaborator = async (data: any): Promise<any> => {
  const response: any = await axiosInstance.post(`/collaborators`, {
    ...data
  });
  return response;
};

export const updateCollaborator = async ({ id, data }: TUpdateCollaborator): Promise<any> => {
  const response: any = await axiosInstance.put(`/collaborators/${id}`, {
    ...data
  });
  return response;
};

export const deleteCollaborator = async (id: string): Promise<any> => {
  const response: any = await axiosInstance.delete(`/collaborators/${id}`);
  return response;
};
