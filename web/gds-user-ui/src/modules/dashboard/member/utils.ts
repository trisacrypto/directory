import { DirectoryType, VaspDirectoryEnum, DirectoryTypeEnum } from './memberType';
import { convertToCVS, downloadCSV } from 'utils/utils';
import { t } from '@lingui/macro';

/* export interface CopyTable {
  label: string;
  value: string;
} */

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
  const m = member.map((item: any) => {
    return {
      ...item,
      registered_directory: getVapsNetwork(item.registered_directory)
    };
  });
  const memberCsv = convertToCVS(m, memberTableHeader as ITableHeader[]);
  downloadCSV(memberCsv, 'members');
};

export async function copyToClipboard(data: any) {
  try {
    const values = data.map((item: any) => {
      item.value = item.value ? item.value : 'N/A';
      return ` ${item.label} \n ${item.value} \n`;
    }).join('\n');
    console.log('text', values);
      await navigator.clipboard.writeText(values);
  } catch (err) {
      console.error('[copyToClipboard]', err);
  }
}
