import { FormControl, FormLabel, Select } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';

const MemberSelectNetwork = () => {
  // this is only responsible for the select network dropdown and set the value to the state
  return (
    <FormControl>
      <FormLabel>
        <Trans>Select Network</Trans>
      </FormLabel>
      <Select>
        <option value="mainnet">MainNet</option>
        <option value="testnet">TestNet</option>
      </Select>
    </FormControl>
  );
};

export default MemberSelectNetwork;
