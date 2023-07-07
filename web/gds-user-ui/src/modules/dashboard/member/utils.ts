import { DirectoryType, VaspDirectoryEnum, DirectoryTypeEnum } from './memberType';
import { convertToCvs, downloadCSV } from 'utils/utils';
import { t } from '@lingui/macro';

export const memberTableHeader = [
  t`Member Name`,
  t`Joined`,
  t`Last Updated`,
  t`Network`,
  t`Status`
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

export const memberToCsv = (member: any) => {
  const memberCsv = convertToCvs(member, memberTableHeader);
  downloadCSV(memberCsv, 'member.csv');
};
