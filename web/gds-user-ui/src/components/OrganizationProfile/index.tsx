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
  TabPanel,
  SimpleGrid,
  List,
  ListItem,
  Table,
  Tbody,
  Tr,
  Td,
  HStack
} from '@chakra-ui/react';

type OrganizationProfileProps = {
  data: any;
  status?: string;
};
const OrganizationProfile: React.FC<OrganizationProfileProps> = (props) => {
  return (
    <Stack py={5} w="full">
      <Stack bg={'#E5EDF1'} h="55px" justifyItems={'center'} p={4} my={5}>
        <Stack mb={10}>
          <Heading fontSize={20}>
            TRISA Organization Profile:{' '}
            <Text as={'span'} color={'blue.500'}>
              [pending registration]
            </Text>
          </Heading>
        </Stack>
      </Stack>
      <HStack py={5}>
        <Stack border={'1px solid #eee'} p={4} my={5} w={'100%'}>
          <Box bg={'white'}>
            <SimpleGrid minChildWidth="120px" spacing="40px">
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
            </SimpleGrid>
          </Box>
        </Stack>
        <Stack border={'1px solid #eee'} p={4} my={5}>
          <Box bg={'white'}>
            <SimpleGrid minChildWidth="120px" spacing="40px">
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
              <List>
                <ListItem>Name Identifier</ListItem>
                <ListItem>ddd</ListItem>
              </List>
            </SimpleGrid>
          </Box>
        </Stack>
      </HStack>
    </Stack>
  );
};

export default OrganizationProfile;
