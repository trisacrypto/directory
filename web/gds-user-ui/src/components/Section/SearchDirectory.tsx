/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { FormEvent, useState, useRef } from 'react';
import {
  Stack,
  Container,
  Box,
  Flex,
  Text,
  Link,
  Button,
  Tooltip,
  FormHelperText,
  FormControl,
  useColorModeValue,
  Table,
  Tr,
  Td,
  Heading,
  HStack,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  TableContainer,
  // Input,
  Tbody
} from '@chakra-ui/react';

import {
  AutoComplete,
  AutoCompleteInput,
  AutoCompleteItem,
  AutoCompleteList
} from '@choc-ui/chakra-autocomplete';

import { colors } from 'utils/theme';
import ErrorMessage from 'components/ui/ErrorMessage';
import countryCodeEmoji, { getCountryName } from 'utils/country';
import { IsoCountryCode } from 'types/type';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';

type TSearchDirectory = {
  handleSubmit: (e: FormEvent, query: string) => void;
  isLoading: boolean;
  result: any;
  error: string;
  query: string;
  handleClose?: () => void;
  onResetData: () => void;
  options: any[];
};
const SearchDirectory: React.FC<TSearchDirectory> = ({
  handleSubmit,
  result,
  error,
  handleClose,
  options,
  onResetData
}) => {
  const [search, setSearch] = useState<string>('');
  const formRef = useRef<HTMLFormElement>(null);

  return (
    <Flex
      py={12}
      width="100%"
      position={'relative'}
      fontFamily={colors.font}
      color={useColorModeValue('black', 'white')}
      minH={400}>
      <Container maxW={'5xl'} zIndex={10} fontFamily={colors.font} mb={10} id={'search'}>
        <Stack>
          <Box mb={{ base: 5 }} color={useColorModeValue('black', 'white')}>
            <Heading fontWeight={600} pb={4} fontSize={'2xl'}>
              <Trans id="Search the Directory Service">Search the Directory Service</Trans>
            </Heading>
            <Text fontSize={{ base: '16px', md: '17px' }}>
              <Trans id="Enter the VASP Common Name or VASP ID. Not a TRISA Member?">
                Enter the VASP Common Name or VASP ID. Not a TRISA Member?
              </Trans>
              <Link href={'/guide'} color={'#1F4CED'} pl={2}>
                <Trans id="Join the TRISA network today.">Join the TRISA network today.</Trans>
              </Link>
            </Text>
          </Box>

          <Stack direction={['column', 'row']} w={'100%'} pb={10}>
            <Text fontSize={'lg'} color={'black'} fontWeight={'semibold'} pt={1}>
              <Trans id="Directory Lookup">Directory Lookup</Trans>
            </Text>
            <Box width={{ md: '70%', sm: '90%' }}>
              <form>
                <FormControl color={'gray.500'}>
                  <HStack>
                    <AutoComplete rollNavigation ref={formRef}>
                      <AutoCompleteInput
                        variant="outline"
                        placeholder="Common name or VASP ID"
                      />

                      <AutoCompleteList>
                        {Object.keys(options)?.map((oid: any, id: any) => (
                          <AutoCompleteItem
                            key={`option-${id}`}
                            value={oid}
                            label={oid}
                            onClick={(e: any) => {
                              onResetData();
                              setSearch(oid);
                              handleSubmit(e, oid);
                            }}
                            >
                            {oid}
                          </AutoCompleteItem>
                        ))}
                      </AutoCompleteList>
                    </AutoComplete>

                    <Button
                      variant="outline"
                      onClick={(e: FormEvent) => {
                        onResetData();
                        formRef.current?.removeItem(search);
                      }}
                      spinnerPlacement="start">
                      Clear
                    </Button>
                  </HStack>

                  <FormHelperText ml={1} color={'#1F4CED'} cursor={'help'}>
                    <Tooltip
                      p={2}
                      label={
                        <>
                          <Text
                            id="TRISA Endpoint is a server address (e.g. trisa.myvasp.com:443) at which
                            the VASP can be reached via secure channels.">
                            TRISA Endpoint is a server address (e.g. trisa.myvasp.com:443) at which
                            the VASP can be reached via secure channels.
                          </Text>
                          <Text>
                            <Trans
                              id="The Common Name typically matches the Endpoint, without the port number
                            at the end (e.g. trisa.myvasp.com) and is used to identify the subject
                            in the X.509 certificate.">
                              The Common Name typically matches the Endpoint, without the port
                              number at the end (e.g. trisa.myvasp.com) and is used to identify the
                              subject in the X.509 certificate.
                            </Trans>
                          </Text>
                        </>
                      }>
                      <Text>
                        <Trans id="What’s a Common name or VASP ID?">
                          What’s a Common name or VASP ID?
                        </Trans>
                      </Text>
                    </Tooltip>
                  </FormHelperText>

                  {error && <ErrorMessage message={error} handleClose={handleClose} />}
                </FormControl>
              </form>
            </Box>
          </Stack>
        </Stack>

        {result && (
          <Box
            textAlign="center"
            justifyContent="center"
            justifyItems={'center'}
            alignContent="center"
            fontSize={'md'}
            mx={'auto'}
            w={'100%'}>
            <Box>
              <Tabs colorScheme="blue">
                <TabList border={'1px solid #eee'}>
                  <Tab
                    sx={{ width: '100%' }}
                    _focus={{ outline: 'none' }}
                    _selected={{ bg: colors.system.blue, color: 'white', fontWeight: 'semibold' }}>
                    <Text fontSize={['x-small', 'medium']}>
                      <Trans id="TESTNET DIRECTORY RECORD">TESTNET DIRECTORY RECORD</Trans>
                    </Text>
                  </Tab>
                  <Tab
                    sx={{ width: '100%' }}
                    _focus={{ outline: 'none' }}
                    _selected={{ bg: colors.system.blue, color: 'white', fontWeight: 'semibold' }}>
                    <Text fontSize={['x-small', 'medium']}>
                      <Trans id="MAINNET DIRECTORY RECORD">MAINNET DIRECTORY RECORD</Trans>
                    </Text>
                  </Tab>
                </TabList>
                <TabPanels>
                  <TabPanel p={0} border="1px solid #eee">
                    <TableContainer>
                      <Table
                        variant="simple"
                        sx={{ 'td:first-child': { fontWeight: 'semibold', width: '50%' } }}>
                        <Tbody>
                          <Tr>
                            <Td>
                              <Trans id="Organization Name">Organization Name</Trans>
                            </Td>
                            <Td colSpan={2}>{result?.testnet?.name || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>
                              <Trans id="Common Name">Common Name</Trans>
                            </Td>
                            <Td>{result?.testnet?.common_name || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>
                              <Trans id="TRISA Service Endpoint">TRISA Service Endpoint</Trans>
                            </Td>
                            <Td>{result?.testnet?.endpoint || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>
                              <Trans id="Registered Directory">Registered Directory</Trans>
                            </Td>
                            <Td>{result?.testnet?.registered_directory || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>
                              <Trans id="TRISA Member ID">TRISA Member ID</Trans>
                            </Td>
                            <Td>{result?.testnet?.id || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>
                              <Trans id="Country">Country</Trans>
                            </Td>
                            <Td>
                              {getCountryName(result?.testnet?.country as IsoCountryCode)}
                              {'  '}
                              {countryCodeEmoji(result?.testnet?.country) || 'N/A'}
                            </Td>
                          </Tr>

                          <Tr>
                            <Td>
                              <Trans id="TRISA Verification">TRISA Verification</Trans>
                            </Td>
                            {result?.testnet?.verified_on ? (
                              <Td>
                                {' '}
                                <Trans id="VERIFIED ON">VERIFIED ON</Trans>{' '}
                                {result?.testnet?.verified_on}{' '}
                              </Td>
                            ) : (
                              <Td>N/A</Td>
                            )}
                          </Tr>
                        </Tbody>
                      </Table>
                    </TableContainer>
                  </TabPanel>
                  <TabPanel p={0} border="1px solid #eee">
                    <TableContainer>
                      <Table
                        variant="simple"
                        sx={{
                          'td:first-child': { fontWeight: 'semibold', width: '50%' },
                          'td:last-child': { width: '50%' }
                        }}>
                        <Tbody>
                          <Tr>
                            <Td>
                              <Trans id="Organization Name">Organization Name</Trans>
                            </Td>
                            <Td colSpan={2}>{result?.mainnet?.name || 'N/A'} </Td>
                          </Tr>
                          <Tr>
                            <Td>
                              <Trans id="Common Name">Common Name</Trans>
                            </Td>
                            <Td>{result?.mainnet?.common_name || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>
                              <Trans id="TRISA Service Endpoint">TRISA Service Endpoint</Trans>
                            </Td>
                            <Td>{result?.mainnet?.endpoint || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>
                              <Trans id="Registered Directory"></Trans>
                            </Td>
                            <Td>{result?.mainnet?.registered_directory || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>
                              <Trans id="TRISA Member ID"></Trans>
                            </Td>
                            <Td>{result?.mainnet?.id || 'N/A'}</Td>
                          </Tr>

                          <Tr>
                            <Td>
                              <Trans id="Country">Country</Trans>
                            </Td>
                            <Td>
                              {getCountryName(result?.mainnet?.country as IsoCountryCode)}
                              {'  '}
                              {countryCodeEmoji(result?.mainnet?.country) || 'N/A'}
                            </Td>
                          </Tr>

                          <Tr>
                            <Td>
                              <Trans id="TRISA Verification">TRISA Verification</Trans>
                            </Td>
                            {result?.mainnet?.verified_on ? (
                              <Td>
                                {' '}
                                <Trans id="VERIFIED ON">VERIFIED ON</Trans>{' '}
                                {result?.mainnet?.verified_on}{' '}
                              </Td>
                            ) : (
                              <Td>N/A</Td>
                            )}
                          </Tr>
                        </Tbody>
                      </Table>
                    </TableContainer>
                  </TabPanel>
                </TabPanels>
              </Tabs>
            </Box>
          </Box>
        )}
      </Container>
    </Flex>
  );
};

SearchDirectory.defaultProps = {
  isLoading: false,
  error: '',
  handleSubmit: () => {},
  result: null
};

export default SearchDirectory;
