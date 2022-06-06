import React, { useState, useEffect } from 'react';
import * as Sentry from '@sentry/react';
import {
  Box,
  Heading,
  VStack,
  Flex,
  Input,
  Stack,
  Text,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel
} from '@chakra-ui/react';

type OrganizationProfileProps = {
  data: any;
  status?: string;
};
const OrganizationProfile: React.FC<OrganizationProfileProps> = (props) => {
  return (
    <Stack py={5} w="full">
      <Stack bg={'#E5EDF1'} h="55px" justifyItems={'center'} p={4} my={5}>
        <Stack>
          <Heading fontSize={20}>
            TRISA Organization Profile:{' '}
            <Text as={'span'} color={'blue.500'}>
              [pending registration]
            </Text>
          </Heading>
        </Stack>
      </Stack>
    </Stack>
  );
};

export default OrganizationProfile;
