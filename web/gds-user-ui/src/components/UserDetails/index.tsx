import React from 'react';

import { Box, Heading, VStack, Flex, Text, Stack } from '@chakra-ui/react';
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
            User Details
          </Heading>
          <Text> User ID : {props.userId} </Text>
          <Text> Profile Created : {props.createdDate}</Text>
          <Text> Status: {props.status}</Text>
          <Text> Permission: {props.permissions} </Text>
          <Text> Last Login: {props.lastLogin} </Text>
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
