import axiosInstance from 'utils/axios';

export const GetAllOrganisations = async () => {
  const response = await axiosInstance.get(`/organizations`);
  return response;
};

export const GetOrganisation = async (id: string) => {
  const response = await axiosInstance.get(`/organizations/${id}`);
  return response;
};

export const CreateOrganisation = async (data: any) => {
  const response = await axiosInstance.post(`/organizations`, data);
  return response;
};

export const UpdateOrganisation = async (id: string, data: any) => {
  const response = await axiosInstance.put(`/organizations/${id}`, data);
  return response;
};
