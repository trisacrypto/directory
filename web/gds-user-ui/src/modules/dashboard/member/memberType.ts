import { t } from '@lingui/macro';

export interface MembersQuery {
  getMembers(): void;
  members: any;
  hasMembersFailed: boolean;
  wasMembersFetched: boolean;
  isFetchingMembers: boolean;
  error?: any;
}

export type DirectoryType = 'testnet' | 'mainnet';
export enum DirectoryTypeEnum {
  TESTNET = 'testnet',
  MAINNET = 'mainnet'
}
export enum VaspDirectoryEnum {
  TESTNET = 'testnet.net',
  MAINNET = 'vaspdirectory.net',
  MAINNET_DEV = 'vaspdirectory.dev',
  TESTNET_DEV = 'testnet.dev'
}

export type VaspType = {
  id: string;
  registered_directory: string;
  common_name: string;
  endpoint: string;
  name: string;
  website: string;
  country: string;
  business_category: string;
  vasp_categories: string[];
  verified_on: string;
  status: string;
  first_listed: string;
  last_updated: string;
};
