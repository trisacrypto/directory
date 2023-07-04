import { DirectoryType, VaspDirectoryEnum, DirectoryTypeEnum } from './memberType';

export const getVaspDirectory = (dir: DirectoryType) => {
  const env = process.env.NODE_ENV;
  if (env === 'development') {
    switch (dir) {
      case DirectoryTypeEnum.MAINNET:
        return VaspDirectoryEnum.MAINNET_DEV;
      case DirectoryTypeEnum.TESTNET:
        return VaspDirectoryEnum.TESTNET_DEV;
      default:
        return VaspDirectoryEnum.MAINNET_DEV;
    }
  } else {
    switch (dir) {
      case DirectoryTypeEnum.MAINNET:
        return VaspDirectoryEnum.MAINNET_PROD;
      case DirectoryTypeEnum.TESTNET:
        return VaspDirectoryEnum.TESTNET_PROD;
      default:
        return VaspDirectoryEnum.MAINNET_PROD;
    }
  }
};
