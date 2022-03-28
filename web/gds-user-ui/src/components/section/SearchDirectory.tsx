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
import ErrorMessage from 'components/ErrorMessage';
import { auto } from '@popperjs/core';
type TSearchDirectory = {
  handleSubmit: (e: FormEvent, query: string) => void;
  isLoading: boolean;
  result: any;
  error: string;
};
const SearchDirectory: React.FC<TSearchDirectory> = ({
  handleSubmit,
  isLoading,
  result,
  error
}) => {
  const [search, setSearch] = useState<string>('');
  return (
    <Flex
      width="100%"
      position={'relative'}
      fontFamily={colors.font}
      color={useColorModeValue('black', 'white')}
      id={'search'}
      minH={400}>
      <Container maxW={'5xl'} fontFamily={colors.font} mb={10}>
        <Stack py={5}>
          <Box mb={{ base: 5 }} color={useColorModeValue('black', 'white')}>
            <Text fontWeight={600} mb={5} fontSize={'2xl'} w={'100%'}>
              Search the Directory Service
            </Text>
            <Text fontSize={'lg'}>
              Not a TRISA Member?
              <Link href={'/register'} color={'#1F4CED'} pl={2}>
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
                  {error && <ErrorMessage message={error} />}
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
            border="2px solid #eee"
            fontSize={18}
            width="100%"
            mx={'auto'}
            w={'2xl'}>
            <Box>
              <Table
                variant={'simple'}
                sx={{
                  'td:first-child': { fontWeight: 'bold' },
                  Tr: { borderStyle: 'hidden' }
                }}>
                <Thead bg={'#eee'} width={'100%'}>
                  <Th>Global TRISA Directory Record</Th>
                  <Th></Th>
                </Thead>
                <Tbody>
                  <Tr>
                    <Td>Common Name</Td>
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
                  <Tr>
                    <Td>Country</Td>
                    <Td>{result.country}</Td>
                  </Tr>
                  <Tr>
                    <Td>TRISA Verification</Td>
                    <Td></Td>
                  </Tr>
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
