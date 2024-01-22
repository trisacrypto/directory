import axiosInstance from 'utils/axios';
export const getAllOrganizations = async (name: string, page: number, pageSize: number) => {
  const urlParams =
    name && name.length > 0
      ? `?name=${encodeURIComponent(name)}&page=${page}&page_size=${pageSize}`
      : `?page=${page}&page_size=${pageSize}`;
  // format the url params

  const response = await axiosInstance.get(`/organizations${urlParams}`);
  return response;
};

export const getOrganization = async (id: string) => {
  const response = await axiosInstance.get(`/organizations/${id}`);
  return response;
};

export const createOrganization = async (data: any) => {
  const response = await axiosInstance.post(`/organizations`, data);
  return response;
};

export const updateOrganization = async (id: string, data: any) => {
  const response = await axiosInstance.put(`/organizations/${id}`, data);
  return response;
};

export const getOrganizationByName = async (name: string, page = 1, pageSize = 8) => {
  const response = await axiosInstance.get(
    `/organizations?name=${name}&page=${page}&page_size=${pageSize}`
  );
  return response;
};
