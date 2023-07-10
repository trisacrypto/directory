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

export const convertMemberToCSV = (jsonData: any, headers: any) => {
  const contentRows = [];

  // Create the header row.
  contentRows.push(headers.join(','));

  // Loop over the member data and create a row for the content displayed in the member summary modal.
  for(let i = 0; i < jsonData.length; i++) {
    const values = headers.map((header: any) => {
      const fieldValue = jsonData[i][header] !== undefined ? jsonData[i][header] : '';
      const escapedValue = fieldValue.toString().replace(/"/g, '""');
      return `"${escapedValue}"`;
    });
    contentRows.push(values.join(','));
  }
  return contentRows.join('\n');
};

export const downloadMemberToCSV = (member: any) => {
  const memberSummary = convertMemberToCSV(member, memberModalHeader);
  downloadCSV(memberSummary, 'member');
};
