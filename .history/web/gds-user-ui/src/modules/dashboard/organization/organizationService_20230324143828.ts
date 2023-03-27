import axiosInstance from 'utils/axios';
export const getAllOrganisations = async (page?: number) => {
  const currentPage = page || 1;
  const pageSize = 8;
  const response = await axiosInstance.get(
    `/organizations?page=${currentPage}&page_size=${pageSize}`
  );
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
