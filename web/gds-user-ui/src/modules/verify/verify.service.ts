import axiosInstance from 'utils/axios';
interface TParams {
  vaspID: any;
  token: any;
  registered_directory: any;
}
export const verifyService = async (query: TParams) => {
  const response = await axiosInstance.get(
    `/verify?vaspID=${query.vaspID}&token=${encodeURIComponent(query.token)}&registered_directory=${
      query.registered_directory
    }`
  );
  return response.data;
};
