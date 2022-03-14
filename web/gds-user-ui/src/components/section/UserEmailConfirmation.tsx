import {
  Flex,
  Box,
  FormControl,
  FormLabel,
  Input,
  Checkbox,
  Stack,
  Link,
  Button,
  Heading,
  Text,
  useColorModeValue,
  Image,
} from '@chakra-ui/react';

import { colors } from '../../utils/theme';

export default function AuthEmailConfirmation() {
  return (
    <Flex
      align={'center'}
      justify={'center'}
      fontFamily={'open sans'}
      fontWeight={'bold'}
      fontSize={'xl'}
      mt={'10%'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={6} mx={'auto'} maxW={'xl'} py={12} px={6} >
        <Stack >
          <Heading fontSize={'xl'}>Thank you for verifying your email address. </Heading>
          <Text color={useColorModeValue('gray.600', 'white')}>

             Your TRISA account in now active. 
          </Text>
         
        </Stack>
          <Stack spacing={8} direction={['column', 'row']} py='10'>
              <Button
                bg={colors.system.blue}
                color={'white'}
                height={'57px'}
                w={['full', '50%']}
                _hover={{
                  bg: '#10aaed',
                }}
                _focus={{
                 borderColor: 'transparent',
               }}
              >
                Log In
              </Button>
              
            </Stack>
      
      </Stack>
    </Flex>
  );
}