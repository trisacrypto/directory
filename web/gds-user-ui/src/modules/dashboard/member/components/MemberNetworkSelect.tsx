import { FormControl, FormLabel, Select } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import React from 'react';
import useMemberState from '../hooks/useMemberState';

const MemberSelectNetwork = () => {
  const { getNetwork, setNetwork } = useMemberState();
  const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setNetwork(event.target.value);
  };

  return (
    <FormControl>
      <FormLabel>
        <Trans>Select Network</Trans>
      </FormLabel>
      <Select defaultValue={getNetwork()} onChange={handleChange} data-testid="select-network">
        <option value="mainnet">MainNet</option>
        <option value="testnet">TestNet</option>
      </Select>
    </FormControl>
  );
};

export default MemberSelectNetwork;
