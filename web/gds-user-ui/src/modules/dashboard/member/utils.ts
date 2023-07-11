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

const memberDetailTableHeader = [
  {
    key: 'name',
    label: t`Name`,
  },
  {
    key: 'website',
    label: t`Website`,
  },
  {
    key: 'business_category',
    label:  t`Business Category`,
  },
  {
    key: 'vasp_categories',
    label: t`VASP Category`,
  },
  {
    key:'country_of_registration',
    label: t`Country of Registration`,
  },
  {
    key: 'technical_contact',
    label:  t`Technical Contact`,
  },
  {
    key: 'legal_contact',
    label:  t`Compliance / Legal Contact`,
  },
  {
    key: 'administrative_contact',
    label: t`Administrative Contact`,
  },
  {
    key: 'trisa_endpoint',
    label: t`TRISA Endpoint`,
  },
  {
    key: 'common_name',
    label: t`Common Name`,
  },
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
  const memberSummary = convertMemberToCSV(member, memberDetailTableHeader);
  downloadCSV(memberSummary, 'member');
};
