import axiosInstance from 'utils/axios';

const getDirectoryUrl = (dir: string) => {
  return dir === 'mainnet' ? 'vaspdirectory.net' : 'trisatest.net';
};

export const getMembersService = async (directory = 'mainnet') => {
  const response = await axiosInstance.get(
    `/members?registered_directory=${getDirectoryUrl(directory)}`
  );
  return response;
};
