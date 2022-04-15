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
  Th
} from '@chakra-ui/react';

import { SearchIcon } from '@chakra-ui/icons';
import { colors } from 'utils/theme';
import ErrorMessage from 'components/ui/ErrorMessage';

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
  const customName = result?.name ? `${result.name} ${query}` : '';

  return (
    <Flex
      width="100%"
      position={'relative'}
      fontFamily={colors.font}
      color={useColorModeValue('black', 'white')}
      minH={400}>
      <Container maxW={'5xl'} zIndex={10} fontFamily={colors.font} mb={10} id={'search'}>
        <Stack py={5}>
          <Box mb={{ base: 5 }} color={useColorModeValue('black', 'white')}>
            <Text fontWeight={600} pb={4} fontSize={'2xl'}>
              Search the Directory Service
            </Text>
            <Text fontSize={'lg'}>
              Enter the VASP Common Name or VASP ID. Not a TRISA Member?
              <Link href={'/getting-started'} color={'#1F4CED'} pl={2}>
                Join the TRISA network today.
              </Link>
            </Text>
          </Box>

          <Stack direction={['column', 'row']} w={'100%'} pb={10}>
            <Text fontSize={'lg'} color={'black'} fontWeight={'bold'}>
              Directory Search
            </Text>
            <Box width={'70%'}>
              <form onSubmit={(e) => handleSubmit(e, search)}>
                <FormControl color={'gray.500'}>
                  <Input
                    size="md"
                    pr="4.5rem"
                    type={'gray.100'}
                    isRequired
                    placeholder="Common name or VASP ID"
                    name="search"
                    onChange={(event) => setSearch(event.currentTarget.value)}
                  />

                  <FormHelperText ml={1} color={'#1F4CED'}>
                    <Tooltip label="TRISA Endpoint is a server address (e.g. trisa.myvasp.com:443) at which the VASP can be reached via secure channels. The Common Name typically matches the Endpoint, without the port number at the end (e.g. trisa.myvasp.com) and is used to identify the subject in the X.509 certificate.">
                      What’s a Common name or VASP ID?
                    </Tooltip>
                  </FormHelperText>
                  <InputRightElement width="2.5rem" color={'black'}>
                    <Button
                      h="2.5rem"
                      size="sm"
                      isLoading={isLoading}
                      variant="outline"
                      type="submit"
                      spinnerPlacement="start">
                      <SearchIcon />
                    </Button>
                  </InputRightElement>
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
            w={'4xl'}>
            <Box>
              <Table
                variant={'ghost'}
                border={'2px solid #eee'}
                css={{ borderCollapse: 'separate', borderSpacing: '0' }}
                sx={{
                  'td:first-child': { fontWeight: 'bold', maxWidth: '25%' },
                  'td:nth-child(2)': { maxWidth: '75%' },
                  Tr: { borderStyle: 'hidden' }
                }}>
                <Thead bg={'#eee'} height={'50px'}>
                  <Th colSpan={2}>Global TRISA Directory Record</Th>
                </Thead>
                <Tbody>
                  <Tr
                    sx={{
                      'td:first-child': { fontWeight: 'bold', width: '30%' },
                      'td:nth-child(2)': { width: '70%' }
                    }}>
                    <Td colSpan={2}>{customName}</Td>
                  </Tr>
                  <Tr>
                    <Td>
                      {/* <Flex minWidth={'100%'} flexWrap="nowrap">
                        <Text minWidth="100%">Common Name</Text>
                      </Flex> */}
                      Common Name
                    </Td>
                    <Td>{result.common_name}</Td>
                  </Tr>
                  <Tr>
                    <Td>TRISA Service Endpoint</Td>
                    <Td>{result.endpoint}</Td>
                  </Tr>
                  <Tr>
                    <Td>Registered Directory</Td>
                    <Td>{result.registered_directory}</Td>
                  </Tr>
                  <Tr>
                    <Td>TRISA Member ID</Td>
                    <Td>{result.id}</Td>
                  </Tr>
                  {result?.country && (
                    <Tr>
                      <Td>Country</Td>
                      <Td>{result.country}</Td>
                    </Tr>
                  )}
                  {result.verified_on && (
                    <Tr>
                      <Td>TRISA Verification</Td>
                      <Td> VERIFIED ON {result.verified_on}</Td>
                    </Tr>
                  )}
                  {result?.identity_certificate?.signature && (
                    <Tr>
                      <Td>TRISA Identity Signature</Td>
                      <Td>
                        <Flex flexWrap="nowrap" wordBreak={'break-word'}>
                          <Text>{result?.identity_certificate?.signature}</Text>
                        </Flex>
                      </Td>
                    </Tr>
                  )}
                </Tbody>
              </Table>
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
