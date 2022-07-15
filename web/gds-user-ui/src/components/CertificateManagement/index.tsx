import DashboardLayout from '../../layouts/DashboardLayout';
import NeedsAttention from '../NeedsAttention';
import {
  Flex,
  Heading,
  Stack,
  Text,
  chakra,
  Box,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  Table,
  TableCaption,
  Thead,
  Tr,
  Th,
  Tbody,
  Button
} from '@chakra-ui/react';
import FormLayout from '../../layouts/FormLayout';
import StatisticCard from './StatisticCard';
import X509TableRows from './X509TableRows';
import SealingCertificateTableRows from './SealingCertificateTableRows';

type CertificateManagementProps = {};

const STATISTICS = { current: 0, expired: 0, revoked: 0, total: 0 };

function CertificateManagement({}: CertificateManagementProps) {
  return (
    <DashboardLayout>
      <Heading marginBottom="30px">Certificate Management</Heading>
      <NeedsAttention text={'Complete Certficate Registration'} buttonText={'Continue'} />

      <Flex
        border="1px solid #DFE0EB"
        fontFamily={'Open Sans'}
        bg={'white'}
        fontSize={'1rem'}
        p={5}
        mt={5}
        boxShadow={'0px 24px 50px rgba(55, 65, 81, 0.25)'}
        borderRadius={'10px'}>
        <Stack>
          <Heading fontSize={'1.2rem'}>TRISA Certificates</Heading>
          <Text>TRISA issues two types of certificates:</Text>
          <Text>
            (1) <chakra.span fontWeight={700}>X.509 Identity Certificates:</chakra.span> X.509
            Identity certificates are about <chakra.span fontWeight={700}>trust</chakra.span>. They
            prove the identity of the organization and are used to maintain a secure connection
            between TRISA members. X.509 Identity certificates must have an Endpoint and Common Name
            and may have Subject Alternative Names associated with them. TRISA’s X.509 Identity
            certificates are valid for one calendar year. TRISA’s X.509 Identity certificates expire
            after one year so member organizations must request a new X.509 Identity certificate
            upon expiration.
          </Text>
          <Text>
            (2) <chakra.span fontWeight={700}>Sealing Certificates:</chakra.span> Sealing
            certificates are for specific use cases such as{' '}
            <chakra.span fontWeight={700}>long-term data storage</chakra.span>. They are used to
            encrypt individual Secure Envelopes or batches of Secure Envelopes at the transactional
            level. Organizations may have multiple signing keys and multiple signing-key
            certificates for different clients, time periods, or use cases. In a compliance
            information transfer, the transactional sealing certificates are referred to as Envelope
            Seals. While an organization may use their X.509 Identity certificate as a sealing
            certificate, TRISA strongly recommends that transactional sealing certificates are
            different from X.509 Identity certificates for security purposes.
          </Text>
        </Stack>
      </Flex>
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
                          <Button>Request New Identity Certificate</Button>
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
      <Box>
        <FormLayout overflowX={'scroll'}>
          <Table variant="unstyled" css={{ borderCollapse: 'separate', borderSpacing: '0 9px' }}>
            <TableCaption placement="top" textAlign="start" p={0} m={0}>
              <Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'}>
                <Heading fontSize={'1.2rem'}>Sealing Certificates</Heading>
                <Button>Request New Sealing Certificate</Button>
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
              <SealingCertificateTableRows />
            </Tbody>
          </Table>
        </FormLayout>
      </Box>
    </DashboardLayout>
  );
}

export default CertificateManagement;
