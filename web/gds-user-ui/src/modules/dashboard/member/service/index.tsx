import axiosInstance from 'utils/axios';
import { DirectoryType } from '../memberType';

const getDirectoryUrl = (dir: DirectoryType) => {
  return dir === 'mainnet' ? 'vaspdirectory.net' : 'trisatest.net';
};

export const getMembersService = async (directory = 'mainnet') => {
  const response = await axiosInstance.get(
    `/members?registered_directory=${getDirectoryUrl(directory as DirectoryType)}`
  );
  return response;
};
