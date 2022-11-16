import { NotAllowedIcon } from '@chakra-ui/icons';
import { Text, VStack } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import React, { ReactNode } from 'react';

type NoDataProps = {
  label?: ReactNode;
};
function NoData({ label }: NoDataProps) {
  return (
    <VStack w="100%" textAlign="center">
      <NotAllowedIcon fontSize="5rem" color="gray.300" />
      <Text textTransform="capitalize" data-testid="label">
        {label || t`No Data`}
      </Text>
    </VStack>
  );
}

export default NoData;
