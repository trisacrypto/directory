import { createSlice } from '@reduxjs/toolkit';
import { DirectoryTypeEnum, type MemberNetworkType } from './memberType';
export const initialValue: MemberNetworkType = {
  network: DirectoryTypeEnum.MAINNET
};

const MemberSlice: any = createSlice({
  name: 'members',
  initialState: initialValue,
  reducers: {
    // get the current network
    getMemberNetwork: (state: any, {}: any) => {
      return state;
    },
    // set the current network
    setMemberNetwork: (state: any, { payload }: any) => {
      state.network = payload;
    },
    // set the default current network
    setDefaultMemberNetwork: (state: any, {}: any) => {
      state.network = 'mainnet';
    }
  }
});

export const { getMemberNetwork, setMemberNetwork, setDefaultMemberNetwork } = MemberSlice.actions;
export const memberReducer = MemberSlice.reducer;
export const memberSelector = (state: any) => state;
