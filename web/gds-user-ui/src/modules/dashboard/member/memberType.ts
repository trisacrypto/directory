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
  TESTNET_DEV = 'testnet.dev',
  TESTNET_PROD = 'testnet.net',
  MAINNET_DEV = 'vaspdirectory.dev',
  MAINNET_PROD = 'vaspdirectory.net'
}
