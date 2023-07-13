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

export const getVaspStatus = (status: any) => {
  switch (status) {
    case 1:
    case 'NO_VERIFICATION':
      return t`NO VERIFICATION`;
    case 2:
    case 'SUBMITTED':
      return t`SUBMITTED`;
    case 3: 
    case 'EMAIL_VERIFIED':
      return t`EMAIL VERIFIED`;
    case 4: 
    case 'PENDING_REVIEW':
      return t`PENDING REVIEW`;
    case 5: 
    case 'REVIEWED':
      return t`REVIEWED`;
    case 6:
    case 'VERIFIED':
      return t`VERIFIED`;
    case 7:
    case 'REJECTED':
      return t`REJECTED`;
    case 8:
    case 'APPEALED':
      return t`APPEALED`;
    case 9:
    case 'ERRORED':
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

export async function copyToClipboard(data: any) {
  try {
    const values = data?.map((item: any) => {
      item.value = item?.value || 'N/A';
      return ` ${item?.label}\n ${item?.value}\n`;
    }).join('\n');
      await navigator.clipboard.writeText(values);
  } catch (err) {
      console.error('[copyToClipboard]', err);
  }
}

export const getBusinessCategory = (category: any) => {
  switch (category) {
    case 1:
    case 'PRIVATE_ORGANIZATION':
      return t`Private Organization`;
    case 2:
    case 'GOVERNMENT_ENTITY':
      return t`Government Entity`;
    case 3:
    case 'BUSINESS_ENTITY':
      return t`Business Entity`;
    case 4:
    case 'NON_COMMERCIAL_ENTITY':
      return t`Non-Commercial Entity`;
    default:
      return t`Unknown Entity`;
  }
};
