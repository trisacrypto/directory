import React, { FormEvent, useState } from 'react';
import {
  Stack,
  Container,
  Box,
  Flex,
  Text,
  Thead,
  Tbody,
  Link,
  Button,
  Tooltip,
  InputRightElement,
  Input,
  FormHelperText,
  FormControl,
  useColorModeValue,
  Table,
  Tr,
  Td,
  Th,
  Heading,
  HStack,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  chakra,
  TableContainer
} from '@chakra-ui/react';

import { SearchIcon } from '@chakra-ui/icons';
import { colors } from 'utils/theme';
import ErrorMessage from 'components/ui/ErrorMessage';
import countryCodeEmoji, { getCountryName } from 'utils/country';
import { IsoCountryCode } from 'types/type';

type TSearchDirectory = {
  handleSubmit: (e: FormEvent, query: string) => void;
  isLoading: boolean;
  result: any;
  error: string;
  query: string;
  handleClose?: () => void;
};
const SearchDirectory: React.FC<TSearchDirectory> = ({
  handleSubmit,
  isLoading,
  result,
  error,
  query,
  handleClose
}) => {
  const [search, setSearch] = useState<string>('');

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
              Search the Directory Service
            </Heading>
            <Text fontSize={{ base: '16px', md: '17px' }}>
              Enter the VASP Common Name or VASP ID. Not a TRISA Member?
              <Link href={'/getting-started'} color={'#1F4CED'} pl={2}>
                Join the TRISA network today.
              </Link>
            </Text>
          </Box>

          <Stack direction={['column', 'row']} w={'100%'} pb={10}>
            <Text fontSize={'lg'} color={'black'} fontWeight={'semibold'} pt={1}>
              Directory Lookup
            </Text>
            <Box width={{ md: '70%', sm: '90%' }}>
              <form onSubmit={(e) => handleSubmit(e, search)}>
                <FormControl color={'gray.500'}>
                  <HStack>
                    <Input
                      size="md"
                      pr="4.5rem"
                      type={'gray.100'}
                      isRequired
                      placeholder="Common name or VASP ID"
                      name="search"
                      onChange={(event) => setSearch(event.currentTarget.value)}
                    />
                    <Button
                      h="2.5rem"
                      size="sm"
                      isLoading={isLoading}
                      variant="outline"
                      type="submit"
                      spinnerPlacement="start">
                      <SearchIcon />
                    </Button>
                  </HStack>

                  <FormHelperText ml={1} color={'#1F4CED'} cursor={'help'}>
                    <Tooltip label="TRISA Endpoint is a server address (e.g. trisa.myvasp.com:443) at which the VASP can be reached via secure channels. The Common Name typically matches the Endpoint, without the port number at the end (e.g. trisa.myvasp.com) and is used to identify the subject in the X.509 certificate.">
                      What’s a Common name or VASP ID?
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
                    <Text fontSize={['x-small', 'medium']}>TESTNET DIRECTORY RECORD</Text>
                  </Tab>
                  <Tab
                    sx={{ width: '100%' }}
                    _focus={{ outline: 'none' }}
                    _selected={{ bg: colors.system.blue, color: 'white', fontWeight: 'semibold' }}>
                    <Text fontSize={['x-small', 'medium']}>MAINNET DIRECTORY RECORD</Text>
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
                            <Td>Organization Name</Td>
                            <Td colSpan={2}>{result[0]?.name}</Td>
                          </Tr>
                          <Tr>
                            <Td>Common Name</Td>
                            <Td>{result[0]?.common_name}</Td>
                          </Tr>
                          <Tr>
                            <Td>TRISA Service Endpoint</Td>
                            <Td>{result[0]?.endpoint}</Td>
                          </Tr>
                          <Tr>
                            <Td>Registered Directory</Td>
                            <Td>{result[0]?.registered_directory}</Td>
                          </Tr>
                          <Tr>
                            <Td>TRISA Member ID</Td>
                            <Td>{result[0]?.id}</Td>
                          </Tr>
                          <Tr>
                            <Td>Country</Td>
                            <Td>
                              {getCountryName(result[0]?.country as IsoCountryCode)}
                              {'  '}
                              {countryCodeEmoji(result[0]?.country) || 'N/A'}
                            </Td>
                          </Tr>

                          <Tr>
                            <Td>TRISA Verification</Td>
                            {result[0]?.verified_on ? (
                              <Td> VERIFIED ON {result[0]?.verified_on} </Td>
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
                            <Td>Organization Name</Td>
                            <Td colSpan={2}>{result[1]?.name || 'N/A'} </Td>
                          </Tr>
                          <Tr>
                            <Td>Common Name</Td>
                            <Td>{result[1]?.common_name || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>TRISA Service Endpoint</Td>
                            <Td>{result[1]?.endpoint || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>Registered Directory</Td>
                            <Td>{result[1]?.registered_directory || 'N/A'}</Td>
                          </Tr>
                          <Tr>
                            <Td>TRISA Member ID</Td>
                            <Td>{result[1]?.id || 'N/A'}</Td>
                          </Tr>

                          <Tr>
                            <Td>Country</Td>
                            <Td>
                              {getCountryName(result[1]?.country as IsoCountryCode)}
                              {'  '}
                              {countryCodeEmoji(result[1]?.country) || 'N/A'}
                            </Td>
                          </Tr>

                          <Tr>
                            <Td>TRISA Verification</Td>
                            {result[1]?.verified_on ? (
                              <Td> VERIFIED ON {result[1]?.verified_on} </Td>
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
