import axiosInstance from 'utils/axios';

export const logUserInBff = async () => {
  const response = await axiosInstance.post(
    `/users/login`
  );

  return response;
};
export const getUserRoles = async () => {
  const response = await axiosInstance.get(`/users/roles`);

  return response;
};

export const getUserCurrentOrganizationAPI = async () => {
  const response = await axiosInstance.get(`/users/organization`);

  return response;
};
