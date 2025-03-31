export interface MembersQuery {
  getMembers(): void;
  members: any;
  hasMembersFailed: boolean;
  wasMembersFetched: boolean;
  isFetchingMembers: boolean;
  error?: any;
}

export interface MemberQuery {
  getMember(): void;
  member: any;
  hasMemberFailed: boolean;
  wasMemberFetched: boolean;
  isFetchingMember: boolean;
  error?: any;
}

export type DirectoryType = 'testnet' | 'mainnet';
export enum DirectoryTypeEnum {
  TESTNET = 'testnet',
  MAINNET = 'mainnet'
}
export enum VaspDirectoryEnum {
  TESTNET = 'testnet.directory',
  MAINNET = 'trisa.directory',
  MAINNET_DEV = 'vaspdirectory.dev',
  TESTNET_DEV = 'trisatest.dev'
}

export type VaspType = {
  id: string;
  registered_directory: string;
  common_name: string;
  endpoint: string;
  name: string;
  website: string;
  country: string;
  business_category: string | number;
  vasp_categories: string[];
  verified_on: string;
  status: string | number;
  first_listed: string;
  last_updated: string;
};

export type MemberDto = {
  vaspId: string;
  network: string;
};

export type MemberSummary = VaspType;

export type Member = {
  data: {
    summary: VaspType;
    legal_person: any;
    trixo: any;
    contacts: any;
  };
};

export type MemberNetworkType = {
  network: DirectoryType;
};
