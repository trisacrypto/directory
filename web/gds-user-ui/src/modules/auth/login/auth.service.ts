import axiosInstance from 'utils/axios';

export const logUserInBff = async (data?: any) => {
  const response = await axiosInstance.post(
    `/users/login`,
    data,
  );
  return response;
};

export const getUserRoles = async () => {
  const response = await axiosInstance.get(`/users/roles`);
  return response;
};

export const getUserCurrentOrganizationService = async () => {
  const response = await axiosInstance.get(`/users/organization`);
  return response;
};

