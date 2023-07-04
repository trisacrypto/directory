export interface membersQuery {
  getMembers(): void;
  members: any;
  hasMembersFailed: boolean;
  wasMembersFetched: boolean;
  isFetchingMembers: boolean;
  error?: any;
}

export type DirectoryType = 'testnet' | 'mainnet';
