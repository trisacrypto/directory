import { DirectoryType, VaspDirectoryEnum, DirectoryTypeEnum } from './memberType';
import { convertToCVS, downloadCSV } from 'utils/utils';
import { t } from '@lingui/macro';

export const memberTableHeader = [
  {
    key: 'name',
    label: t`Member Name`
  },
  {
    key: 'joined',
    label: t`Joined`
  },
  {
    key: 'last_updated',
    label: t`Last Updated`
  },
  {
    key: 'network',
    label: t`Network`
  },
  {
    key: 'status',
    label: t`Status`
  }
];
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

export const downloadMembers2CVS = (member: any) => {
  const memberCsv = convertToCVS(member, memberTableHeader as ITableHeader[]);
  downloadCSV(memberCsv, 'members');
};
