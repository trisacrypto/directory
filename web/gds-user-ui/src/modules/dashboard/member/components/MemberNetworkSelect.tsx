import { FormControl, FormLabel, Select } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import { Dispatch } from '@reduxjs/toolkit';
import React from 'react';
import { useDispatch } from 'react-redux';
import { setMemberNetwork } from '../member.slice';
import Store from 'application/store';

const MemberSelectNetwork = () => {
  const dispatch: Dispatch = useDispatch();
  const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
  dispatch(setMemberNetwork(event.target.value));
};

// Set the default value of the network to the value in the store.
const { network } = Store.getState().members;

  return (
    <FormControl>
      <FormLabel>
        <Trans>Select Network</Trans>
      </FormLabel>
      <Select defaultValue={network} onChange={handleChange} data-testid="select-network">
        <option value="mainnet">MainNet</option>
        <option value="testnet">TestNet</option>
      </Select>
    </FormControl>
  );
};

export default MemberSelectNetwork;
