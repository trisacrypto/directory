import { DirectoryType, VaspDirectoryEnum, DirectoryTypeEnum } from './memberType';

export const getVaspDirectory = (dir: DirectoryType) => {
  return dir === DirectoryTypeEnum.TESTNET ? VaspDirectoryEnum.TESTNET : VaspDirectoryEnum.MAINNET;
};

export const getVapsNetwork = (dir: any) => {
  switch (dir) {
    case (VaspDirectoryEnum.TESTNET, VaspDirectoryEnum.TESTNET_DEV):
      return 'TestNet';
    case (VaspDirectoryEnum.MAINNET, VaspDirectoryEnum.MAINNET_DEV):
      return 'MainNet';
    default:
      return 'MainNet';
  }
};
