import React from 'react';

import { Box, Heading, VStack, Flex, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
type UserDetailsProps = {
  userId: string;
  createdDate: string;
  status: string;
  permissions: string;
  lastLogin: string;
};
const UserDetails: React.FC<UserDetailsProps> = ({
  userId,
  createdDate,
  permissions,
  lastLogin,
  status
}) => {
  return (
    <Flex mt={10}>
      <VStack spacing={4}>
        <Box mt={2}>
          <Heading pb={3} size="md">
            <Trans id="User Details">User Details</Trans>
          </Heading>
          <Text data-testid="user_id">
            <Trans id="User ID:">User ID:</Trans> {userId}
          </Text>
          <Text data-testid="profile_created">
            <Trans id="Profile Created:">Profile Created:</Trans> {createdDate}
          </Text>
          <Text data-testid="status">
            <Trans id="Status:">Status:</Trans> {status}
          </Text>
          <Text data-testid="permissions">
            <Trans id="Permission:">Permission:</Trans> {permissions}
          </Text>
          <Text data-testid="last_login">
            <Trans id="Last Login:">Last Login:</Trans> {lastLogin}
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
