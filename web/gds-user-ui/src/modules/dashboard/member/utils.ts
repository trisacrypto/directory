import { DirectoryType, VaspDirectoryEnum, DirectoryTypeEnum } from './memberType';

export const getVaspDirectory = (dir: DirectoryType) => {
  return dir === DirectoryTypeEnum.TESTNET ? VaspDirectoryEnum.TESTNET : VaspDirectoryEnum.MAINNET;
};
