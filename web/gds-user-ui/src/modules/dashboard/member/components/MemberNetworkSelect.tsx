import { FormControl, FormLabel, Select } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import { Dispatch } from '@reduxjs/toolkit';
import React from 'react';
import { useDispatch } from 'react-redux';
import { setMemberNetwork } from '../member.slice';

const MemberSelectNetwork = () => {
  const dispatch: Dispatch = useDispatch();
  const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
  dispatch(setMemberNetwork(event.target.value));
};

  return (
    <FormControl>
      <FormLabel>
        <Trans>Select Network</Trans>
      </FormLabel>
      <Select data-testid="select-network" onChange={handleChange}>
        <option value="mainnet">MainNet</option>
        <option value="testnet">TestNet</option>
      </Select>
    </FormControl>
  );
};

export default MemberSelectNetwork;
