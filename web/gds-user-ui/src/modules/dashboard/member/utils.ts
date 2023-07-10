import { DirectoryType, VaspDirectoryEnum, DirectoryTypeEnum } from './memberType';
import { convertToCVS, downloadCSV } from 'utils/utils';
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

export const downloadMembers2CVS = (member: any) => {
  const memberCsv = convertToCVS(member, memberTableHeader);
  downloadCSV(memberCsv, 'members');
};

const memberModalHeader = [
  t`Member Name`,
  t`Website`,
  t`Business Category`,
  t`VASP Category`,
  t` Country of Registration`,
  t`Technical Contact`,
  t`Compliance / Legal Contact`,
  t`Administrative Contact`,
  t`TRISA Endpoint`,
  t`Common Name`
];

export const downloadMemberToCSV = (member: any) => {
  const memberSummary = convertToCVS(member, memberModalHeader);
  downloadCSV(memberSummary, 'member');
};
