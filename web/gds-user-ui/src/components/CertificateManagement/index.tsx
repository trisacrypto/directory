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
  Accordion,
  AccordionButton,
  AccordionItem,
  AccordionPanel,
  AccordionIcon,
  AlertIcon,
  Alert,
  HStack,
  Button,
  AlertDescription
} from '@chakra-ui/react';
import MainnetTestnetCertificates from './MainnetTestnetCertificates';
import { t, Trans } from '@lingui/macro';
import useGetCertificates from 'hooks/useGetCertificates';
import dayjs from 'dayjs';

const CertificateValidityAlert = ({ hasValidMainnet, hasValidTestnet }: any) => (
  <>
    {hasValidMainnet && hasValidTestnet ? (
      <Alert bg="#D8EAF6" borderRadius={'10px'} mb={2}>
        <AlertIcon />
        <HStack justifyContent={'space-between'} w="100%">
          <AlertDescription>
            <Trans>
              The Organization has a current and valid{' '}
              <chakra.span fontWeight={700}>Mainnet</chakra.span> and{' '}
              <chakra.span fontWeight={700}>Tesnet</chakra.span> Identity Certificate.
            </Trans>
          </AlertDescription>
          <Button
            border={'1px solid white'}
            width={142}
            px={8}
            as={'a'}
            borderRadius={0}
            color="#fff"
            cursor="pointer"
            bg="#000"
            _hover={{ bg: '#000000D1' }}>
            <Trans>View/Edit</Trans>
          </Button>
        </HStack>
      </Alert>
    ) : null}

    {hasValidMainnet ? (
      <Alert bg="#D8EAF6" borderRadius={'10px'} mb={2}>
        <AlertIcon />
        <HStack justifyContent={'space-between'} w="100%">
          <AlertDescription>
            <Trans>
              The Organization has a current and valid{' '}
              <chakra.span fontWeight={700}>Mainnet</chakra.span> Identity Certificate.
            </Trans>
          </AlertDescription>
          <Button
            border={'1px solid white'}
            width={142}
            px={8}
            as={'a'}
            borderRadius={0}
            color="#fff"
            cursor="pointer"
            bg="#000"
            _hover={{ bg: '#000000D1' }}>
            <Trans>View/Edit</Trans>
          </Button>
        </HStack>
      </Alert>
    ) : null}

    {hasValidTestnet ? (
      <Alert bg="#D8EAF6" borderRadius={'10px'} mb={2}>
        <AlertIcon />
        <HStack justifyContent={'space-between'} w="100%">
          <AlertDescription>
            <Trans>
              The Organization has a current and valid{' '}
              <chakra.span fontWeight={700}>Testnet</chakra.span> Identity Certificate.
            </Trans>
          </AlertDescription>
          <Button
            border={'1px solid white'}
            width={142}
            px={8}
            as={'a'}
            borderRadius={0}
            color="#fff"
            cursor="pointer"
            bg="#000"
            _hover={{ bg: '#000000D1' }}>
            <Trans>View/Edit</Trans>
          </Button>
        </HStack>
      </Alert>
    ) : null}
  </>
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const hasExpired = (date: string) => {
  const _date = dayjs(date);
  if (_date.isValid()) {
    return dayjs(_date).isBefore(dayjs());
  }

  return false;
};

function CertificateManagement() {
  const { data: certificates } = useGetCertificates();

  return (
    <>
      <Heading marginBottom="30px">
        <Trans id="Certificate Management">Certificate Inventory</Trans>
      </Heading>
      <CertificateValidityAlert hasValidMainnet={true} hasValidTestNet={false} />
      <NeedsAttention text={t`Complete Certficate Registration`} buttonText={t`Continue`} />
      <Flex
        border="1px solid #DFE0EB"
        fontFamily={'Open Sans'}
        bg={'white'}
        fontSize={'1rem'}
        px={5}
        py={3}
        mt={5}
        borderRadius={'10px'}>
        <Stack w="100%">
          <Accordion allowToggle>
            <AccordionItem border="none">
              <AccordionButton
                display="flex"
                justifyContent="space-between"
                px={0}
                _hover={{ bg: 'none' }}>
                <Heading fontSize={'1.2rem'}>
                  <Trans>TRISA Certificate Types</Trans>
                </Heading>
                <AccordionIcon />
              </AccordionButton>
              <AccordionPanel px={0}>
                <Text>
                  <Trans id="TRISA issues two types of certificates:">
                    TRISA issues two types of certificates:
                  </Trans>
                </Text>
                <Text mt={2}>
                  <Trans>
                    (1){' '}
                    <chakra.span fontWeight={700}>
                      <Trans>X.509 Identity Certificates:</Trans>
                    </chakra.span>{' '}
                    X.509 Identity certificates are about{' '}
                    <chakra.span fontWeight={700}>trust</chakra.span>. They prove the identity of
                    the organization and are used to maintain a secure connection between TRISA
                    members. X.509 Identity certificates must have an Endpoint and Common Name and
                    may have Subject Alternative Names associated with them. TRISA’s X.509 Identity
                    certificates are valid for one calendar year. TRISA’s X.509 Identity
                    certificates expire after one year so member organizations must request a new
                    X.509 Identity certificate upon expiration.
                  </Trans>
                </Text>
                <Text mt={3}>
                  <Trans>
                    (2){' '}
                    <chakra.span fontWeight={700}>
                      <Trans>Sealing Certificates:</Trans>
                    </chakra.span>{' '}
                    Sealing certificates are for specific use cases such as{' '}
                    <chakra.span fontWeight={700}>long-term data storage</chakra.span>. They are
                    used to encrypt individual Secure Envelopes or batches of Secure Envelopes at
                    the transactional level. Organizations may have multiple signing keys and
                    multiple signing-key certificates for different clients, time periods, or use
                    cases. In a compliance information transfer, the transactional sealing
                    certificates are referred to as Envelope Seals. While an organization may use
                    their X.509 Identity certificate as a sealing certificate, TRISA strongly
                    recommends that transactional sealing certificates are different from X.509
                    Identity certificates for security purposes.
                  </Trans>
                </Text>
              </AccordionPanel>
            </AccordionItem>
          </Accordion>
        </Stack>
      </Flex>
      <Heading fontSize={'1.2rem'} fontWeight={700} p={5} my={5} bg={'#E5EDF1'}>
        <Trans>X.509 Identity Certificate Inventory</Trans>
      </Heading>
      <Box>
        <Tabs isFitted>
          <TabList bg={'#E5EDF1'} border={'1px solid rgba(0, 0, 0, 0.29)'}>
            <Tab _selected={{ bg: '#60C4CA', fontWeight: 700 }}>
              <Trans>MainNet Certificates</Trans>
            </Tab>
            <Tab _selected={{ bg: '#60C4CA', fontWeight: 700 }}>
              <Trans>TestNet Certificates</Trans>
            </Tab>
          </TabList>
          <TabPanels>
            <TabPanel>
              <MainnetTestnetCertificates network="mainnet" data={certificates?.mainnet} />
            </TabPanel>
            <TabPanel>
              <MainnetTestnetCertificates network="testnet" data={certificates?.testnet} />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>
      {/* <Heading fontSize={'1.2rem'} fontWeight={700} p={5} my={5} bg={'#E5EDF1'} mx={4}>
        <Trans>Sealing Certificate Inventory</Trans>
      </Heading>
      <Box px={4}>
        <CertificateDataTable />
      </Box> */}
    </>
  );
}

export default CertificateManagement;
