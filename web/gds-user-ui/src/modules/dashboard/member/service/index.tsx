import axiosInstance from 'utils/axios';
import { DirectoryType, MemberDto } from '../memberType';
import { getVaspDirectory } from '../utils';

export const getMembersService = async (directory: DirectoryType) => {
  const response = await axiosInstance.get(
    `/members?registered_directory=${getVaspDirectory(directory)}`
  );
  return response;
};

export const getMemberService = async (payload: MemberDto) => {
  const { vaspId, network } = payload;
  const url = `/members/${vaspId}?registered_directory=${getVaspDirectory(
    network as DirectoryType
  )}`;
  const response = await axiosInstance.get(url);

  return response;
};
