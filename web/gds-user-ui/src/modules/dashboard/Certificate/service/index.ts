import axiosInstance from 'utils/axios';

export const lookup = async (params: string, body: any) => {
  const response = await axiosInstance.post(`/register/${params}`, { ...body });
  return response.data;
};
