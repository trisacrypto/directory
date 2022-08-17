import React from 'react';

import { Box, Heading, VStack, Flex, Text, Stack } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
type UserDetailsProps = {
  userId: string;
  createdDate: string;
  status: string;
  permissions: string;
  lastLogin: string;
};
const UserDetails: React.FC<UserDetailsProps> = (props) => {
  return (
    <Flex mt={10}>
      <VStack spacing={4}>
        <Box mt={2}>
          <Heading pb={3} size="md">
            <Trans id="User Details">User Details</Trans>
          </Heading>
          <Text>
            <Trans id="User ID:">User ID:</Trans> {props.userId}
          </Text>
          <Text>
            <Trans id="Profile Created:">Profile Created:</Trans> {props.createdDate}
          </Text>
          <Text>
            <Trans id="Status:">Status:</Trans> {props.status}
          </Text>
          <Text>
            <Trans id="Permission:">Permission:</Trans> {props.permissions}
          </Text>
          <Text>
            <Trans id="Last Login:">Last Login:</Trans> {props.lastLogin}
          </Text>
        </Box>
      </VStack>
    </Flex>
  );
};

UserDetails.defaultProps = {
  userId: 'C0000213',
  createdDate: '01/01/2020',
  status: 'Active',
  permissions: 'Admin',
  lastLogin: '01/01/2020'
};

export default UserDetails;
