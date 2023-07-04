import axiosInstance from 'utils/axios';
import { DirectoryType } from '../memberType';
import { getVaspDirectory } from '../utils';

export const getMembersService = async (directory = 'mainnet') => {
  const response = await axiosInstance.get(
    `/members?registered_directory=${getVaspDirectory(directory as DirectoryType)}`
  );
  return response;
};
