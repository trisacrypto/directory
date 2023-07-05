import { DirectoryType, VaspDirectoryEnum, DirectoryTypeEnum } from './memberType';

export const getVaspDirectory = (dir: DirectoryType) => {
  return dir === DirectoryTypeEnum.TESTNET ? VaspDirectoryEnum.TESTNET : VaspDirectoryEnum.MAINNET;
};

export const getVapsNetwork = (dir: any) => {
  switch (dir) {
    case (VaspDirectoryEnum.TESTNET, VaspDirectoryEnum.TESTNET_DEV):
      return VaspDirectoryEnum.TESTNET;
    case (VaspDirectoryEnum.MAINNET, VaspDirectoryEnum.MAINNET_DEV):
      return VaspDirectoryEnum.MAINNET;
    default:
      return VaspDirectoryEnum.TESTNET;
  }
};
