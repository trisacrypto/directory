import React from "react";
import {
  Stack,
  Container,
  Box,
  Flex,
  Text,
  Heading,
  Button,
  Tooltip,
  InputRightElement,
  Input,
  FormHelperText,
  FormControl,
  useColorModeValue
} from "@chakra-ui/react";

import { SearchIcon } from '@chakra-ui/icons'

import { colors } from '../../utils/theme';

export default function SearchDirectory() {
    return (
      
 
       <Flex
            minH={'100vh'}
            minWidth={'100vw'}
            bg={useColorModeValue('white', 'gray.800')}
            position={"relative"}
            color={useColorModeValue('black', 'white')} >
      <Container maxW={"3xl"} zIndex={10} position={"relative"}>
        <Stack>
          <Stack
            flex={1}
            justify={{ lg: "center" }}
            py={{ base: 4, md: 20 }}
            px={{ base: 10, md: 55 }}
          >
                      <Box mb={{ base: 2, md: 20 }} color={useColorModeValue('blak','white')}>
              <Heading
                fontFamily={"heading"}
                mb={3}
                fontSize={"xl"}
              >
               Search the Directory Service
              </Heading>
              <Text fontSize={"lg"} >
              Not a TRISA Member? Join the TRISA network today.
              </Text>
            </Box>
        
                <Stack  direction={['column','row']}>
                  
                  <Text
                    fontFamily={"Open Sans"}
                    fontSize={"lg"}
                    color={"black"}
                  >
                    Directory Search
                  </Text>
                  
                   <FormControl color={'gray.500'} >
                   
                   <Input size='md'
                    pr='4.5rem'
                    type={'gray.100'}
                    placeholder='Common name or VASP ID'
                   />
                   
                   <FormHelperText ml={1} color={'#1F4CED'}>
                     <Tooltip label=' Some description about name and vasp id'>
                    What’s a Common name or VASP ID?
                    </Tooltip>
                   </FormHelperText>
                   <InputRightElement width='2.5rem' color={'black'}>
                    <Button h='2.5rem' size='sm' onClick={e=>{}}>
                      <SearchIcon/>
                    </Button>
                  </InputRightElement>
                  </FormControl>
             
                </Stack>
               
             <Box alignItems="center" pt={10} textAlign="center">
             
            </Box>
          </Stack>
          <Flex flex={1} />
        </Stack>
      </Container>
      </Flex>
   
  );
}
