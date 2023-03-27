import axiosInstance from 'utils/axios';
import type { OrganizationPagination } from './organizationType';
export const GetAllOrganisations = async (params: OrganizationPagination) => {
  const page = params.page || 1;
  const pageSize = params.pageSize || 8;
  const response = await axiosInstance.get(`/organizations?page=${page}&page_size=${pageSize}`);
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
