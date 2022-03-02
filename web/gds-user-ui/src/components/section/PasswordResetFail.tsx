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

import { GoogleIcon } from '../icon';

import { colors } from '../../utils/theme';

export default function CreateAccount() {
  return (
    <Flex
      align={'center'}
      justify={'center'}
      fontFamily={'open sans'}
      fontSize={'xl'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={12} mx={'auto'} maxW={'lg'} py={12} px={6} >
        <Stack align={'center'}>
          <Heading fontSize={'xl'}>Sorry. We could not find a user account with the email address [insert email address].  </Heading>
          
         
        </Stack>
      
          <Box
          rounded={'lg'}
          bg={useColorModeValue('white', 'transparent')}
          position={'relative'}
          bottom={5}
         >
          <Text textAlign='center'>
              Not a TRISA Member? <Link href="/register" color={colors.system.cyan}> Join the TRISA network today. </Link>
              </Text>
        </Box>
      </Stack>
    </Flex>
  );
}