import {
<<<<<<< HEAD
    Box,
    Heading,
    Stack,
    Tab,
    Table,
    TableCaption,
    TabList,
    TabPanel,
    TabPanels,
    Tabs, Tbody, Th,
    Thead, Tr
} from "@chakra-ui/react";
import FormLayout from "../../layouts/FormLayout";
import FormButton from "../ui/FormButton";
import React from "react";
import StatisticCard from "./StatisticCard";
import X509TableRows from "./X509TableRows";

const STATISTICS = { current: 0, expired: 0, revoked: 0, total: 0 };

function X509IdentityCertificateInventory(){

    return (
        <>
            <Heading fontSize={'1.2rem'} fontWeight={700} p={5} my={5} bg={"#E5EDF1"}>X.509 Identity Certificate Inventory</Heading>
            <Box>
                <Tabs isFitted>
                    <TabList bg={"#E5EDF1"} border={"1px solid rgba(0, 0, 0, 0.29)"}>
                        <Tab _selected={{ bg: "#60C4CA", fontWeight: 700 }}>MainNet Certificates</Tab>
                        <Tab _selected={{ bg: "#60C4CA", fontWeight: 700 }}>TestNet Certificates</Tab>
                    </TabList>
                    <TabPanels>
                        <TabPanel>
                            <Stack spacing={5}>
                                <Stack direction={"row"} flexWrap={'wrap'} spacing={4}>
                                    {
                                        Object.entries(STATISTICS).map(statistic => (
                                            <StatisticCard key={statistic[0]} title={statistic[0]} total={statistic[1]} />
                                        ))
                                    }
                                </Stack>
                                <Box>
                                    <FormLayout overflowX={'scroll'}>
                                        <Table variant="unstyled" css={{ borderCollapse: 'separate', borderSpacing: '0 9px' }}>
                                            <TableCaption placement="top" textAlign="start" p={0} m={0} >
                                                <Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'}>
                                                    <Heading fontSize={'1.2rem'}>X.509 Identity Certificates</Heading>
                                                    <FormButton borderRadius={5}>Request New Identity Certificate</FormButton>
                                                </Stack>
                                            </TableCaption>
                                            <Thead>
                                                <Tr>
                                                    <Th>No</Th>
                                                    <Th>Signature ID</Th>
                                                    <Th>Issue Date</Th>
                                                    <Th>Expiration Date</Th>
                                                    <Th>Status</Th>
                                                    <Th textAlign="center">Action</Th>
                                                </Tr>
                                            </Thead>
                                            <Tbody>
                                                <X509TableRows />
                                            </Tbody>
                                        </Table>
                                    </FormLayout>
                                </Box>
                            </Stack>
                        </TabPanel>
                        <TabPanel>
                            <p>two!</p>
                        </TabPanel>
                    </TabPanels>
                </Tabs>
            </Box>
            <Heading fontSize={'1.2rem'} fontWeight={700} p={5} my={5} bg={"#E5EDF1"}>Sealing Certificate Inventory</Heading>
        </>
    )
}

export default X509IdentityCertificateInventory
=======
  Box,
  Button,
  Heading,
  Stack,
  Tab,
  Table,
  TableCaption,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Tbody,
  Th,
  Thead,
  Tr
} from '@chakra-ui/react';
import FormLayout from '../../layouts/FormLayout';
import StatisticCard from './StatisticCard';
import X509TableRows from './X509TableRows';

const STATISTICS = { current: 0, expired: 0, revoked: 0, total: 0 };

function X509IdentityCertificateInventory() {
  return (
    <>
      <Heading fontSize={'1.2rem'} fontWeight={700} p={5} my={5} bg={'#E5EDF1'}>
        X.509 Identity Certificate Inventory
      </Heading>
      <Box>
        <Tabs isFitted>
          <TabList bg={'#E5EDF1'} border={'1px solid rgba(0, 0, 0, 0.29)'}>
            <Tab _selected={{ bg: '#60C4CA', fontWeight: 700 }}>MainNet Certificates</Tab>
            <Tab _selected={{ bg: '#60C4CA', fontWeight: 700 }}>TestNet Certificates</Tab>
          </TabList>
          <TabPanels>
            <TabPanel>
              <Stack spacing={5}>
                <Stack direction={'row'} flexWrap={'wrap'} spacing={4}>
                  {Object.entries(STATISTICS).map((statistic) => (
                    <StatisticCard key={statistic[0]} title={statistic[0]} total={statistic[1]} />
                  ))}
                </Stack>
                <Box>
                  <FormLayout overflowX={'scroll'}>
                    <Table
                      variant="unstyled"
                      css={{ borderCollapse: 'separate', borderSpacing: '0 9px' }}>
                      <TableCaption placement="top" textAlign="start" p={0} m={0}>
                        <Stack
                          direction={'row'}
                          alignItems={'center'}
                          justifyContent={'space-between'}>
                          <Heading fontSize={'1.2rem'}>X.509 Identity Certificates</Heading>
                          <Button borderRadius={5}>Request New Identity Certificate</Button>
                        </Stack>
                      </TableCaption>
                      <Thead>
                        <Tr>
                          <Th>No</Th>
                          <Th>Signature ID</Th>
                          <Th>Issue Date</Th>
                          <Th>Expiration Date</Th>
                          <Th>Status</Th>
                          <Th textAlign="center">Action</Th>
                        </Tr>
                      </Thead>
                      <Tbody>
                        <X509TableRows />
                      </Tbody>
                    </Table>
                  </FormLayout>
                </Box>
              </Stack>
            </TabPanel>
            <TabPanel>
              <p>two!</p>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>
      <Heading fontSize={'1.2rem'} fontWeight={700} p={5} my={5} bg={'#E5EDF1'}>
        Sealing Certificate Inventory
      </Heading>
    </>
  );
}

export default X509IdentityCertificateInventory;
>>>>>>> origin/main
