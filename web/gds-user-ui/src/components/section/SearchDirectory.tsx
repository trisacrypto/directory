import React from 'react';
import {
  Stack,
  Container,
  Box,
  Flex,
  Text,
  Heading,
  Link,
  Button,
  Tooltip,
  InputRightElement,
  Input,
  FormHelperText,
  FormControl,
  useColorModeValue
} from '@chakra-ui/react';

import { SearchIcon } from '@chakra-ui/icons';
import { colors } from 'utils/theme';
export default function SearchDirectory() {
  return (
    <Box
      width="100%"
      position={'relative'}
      fontFamily={colors.font}
      color={useColorModeValue('black', 'white')}
      id={'search'}>
      <Container maxW={'5xl'} zIndex={10} position={'relative'} fontFamily={colors.font} mb={10}>
        <Stack>
          <Stack flex={1} justify={{ lg: 'center' }} py={{ base: 4, md: 10 }}>
            <Box mb={{ base: 5 }} color={useColorModeValue('black', 'white')}>
              <Text fontWeight={600} mb={5} fontSize={'2xl'}>
                Search the Directory Service
              </Text>
              <Text fontSize={'lg'}>
                Not a TRISA Member?
                <Link href={'/register'} color={'#1F4CED'} pl={2}>
                  Join the TRISA network today.
                </Link>
              </Text>
            </Box>

            <Stack direction={['column', 'row']}>
              <Text fontSize={'lg'} color={'black'}>
                Directory Search
              </Text>

              <FormControl color={'gray.500'}>
                <Input
                  size="md"
                  pr="4.5rem"
                  type={'gray.100'}
                  placeholder="Common name or VASP ID"
                />

                <FormHelperText ml={1} color={'#1F4CED'}>
                  <Tooltip label="TRISA Endpoint is a server address (e.g. trisa.myvasp.com:443) at which the VASP can be reached via secure channels. The Common Name typically matches the Endpoint, without the port number at the end (e.g. trisa.myvasp.com) and is used to identify the subject in the X.509 certificate.">
                    What’s a Common name or VASP ID?
                  </Tooltip>
                </FormHelperText>
                <InputRightElement width="2.5rem" color={'black'}>
                  <Button h="2.5rem" size="sm" onClick={(e) => {}}>
                    <SearchIcon />
                  </Button>
                </InputRightElement>
              </FormControl>
            </Stack>

            <Box alignItems="center" pt={10} textAlign="center"></Box>
          </Stack>
          <Flex flex={1} />
        </Stack>
      </Container>
    </Box>
  );
}
