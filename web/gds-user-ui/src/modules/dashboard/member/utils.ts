import { DirectoryType, VaspDirectoryEnum, DirectoryTypeEnum } from './memberType';
import { convertToCVS, downloadCSV } from 'utils/utils';
import { t } from '@lingui/macro';

export const memberTableHeader = [
  {
    key: 'name',
    label: t`Member Name`
  },
  {
    key: 'first_listed',
    label: t`Joined`
  },
  {
    key: 'last_updated',
    label: t`Last Updated`
  },
  {
    key: 'registered_directory',
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

export const getVaspNetwork = (dir: any) => {
  switch (dir) {
    case (VaspDirectoryEnum.TESTNET, VaspDirectoryEnum.TESTNET_DEV):
      return 'TestNet';
    case (VaspDirectoryEnum.MAINNET, VaspDirectoryEnum.MAINNET_DEV):
      return 'MainNet';
    default:
      return 'MainNet';
  }
};

export const getVaspStatus = (status: number) => {
  switch (status) {
    case 1:
      return t`NO VERIFICATION`;
    case 2:
      return t`SUBMITTED`;
    case 3:
      return t`EMAIL VERIFIED`;
    case 4:
      return t`PENDING REVIEW`;
    case 5:
      return t`REVIEWED`;
    case 6:
      return t`VERIFIED`;
    case 7:
      return t`REJECTED`;
    case 8:
      return t`APPEALED`;
    case 9:
      return t`ERRORED`;
    default:
      return t`NO VERIFICATION`;
  }
};

export const downloadMembers2CVS = (member: any) => {
  const m = member.map((item: any) => {
    return {
      ...item,
      status: getVaspStatus(item.status),
      registered_directory: getVaspNetwork(item.registered_directory)
    };
  });
  const memberCsv = convertToCVS(m, memberTableHeader as ITableHeader[]);
  downloadCSV(memberCsv, 'members');
};

