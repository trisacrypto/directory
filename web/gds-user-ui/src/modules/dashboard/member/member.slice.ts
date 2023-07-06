import { createSlice } from '@reduxjs/toolkit';
import type { MemberNetworkType } from './memberType';
export const initialValue: MemberNetworkType = {
  network: 'mainnet'
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
    }
  }
});

export const { getMemberNetwork, setMemberNetwork } = MemberSlice.actions;
export const memberReducer = MemberSlice.reducer;
export const memberSelector = (state: any) => state;
