import axiosInstance from 'utils/axios';

export const getMembersService = async () => {
  const response = await axiosInstance.get(`/members`);
  return response;
};
