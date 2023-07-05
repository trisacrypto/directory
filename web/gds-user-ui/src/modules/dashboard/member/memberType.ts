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
  TESTNET = 'testnet.net',
  MAINNET = 'vaspdirectory.net'
}
