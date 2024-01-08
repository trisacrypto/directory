import { Heading, Stack, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import React from 'react';
import AddressList from './AddressList';

const Addresses: React.FC = () => {
  return (
    <Stack pt={5}>
      <Heading size="md">
        <Trans id="Addresses">Addresses</Trans>
      </Heading>
      <Text size="sm">
        <Trans id="At least one geographic address is required. Enter the primary geographic address of the organization. Organizations may enter additional addresses if operating in multiple jurisdictions.">
          At least one geographic address is required. Enter the primary geographic address of the
          organization. Organizations may enter additional addresses if operating in multiple
          jurisdictions.
        </Trans>
      </Text>
      <AddressList />
    </Stack>
  );
};

export default Addresses;
