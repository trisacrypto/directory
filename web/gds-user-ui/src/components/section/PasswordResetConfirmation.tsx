import {
  Flex,
  Box,
  Stack,
  Link,
  Heading,
  Text,
  useColorModeValue,
 
} from '@chakra-ui/react';

import { colors } from '../../utils/theme';

export default function PasswordResetConfirmation(props : any) {
  return (
    <Flex
      align={'center'}
      justify={'center'}
      fontFamily={'open sans'}
      fontSize={'xl'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={12} mx={'auto'} maxW={'lg'} py={12} px={6} >
        <Stack align={'center'}>
          <Heading fontSize={'xl'}>
            Thank you. We have sent instructions to reset your password to
             <Text as={'span'} fontWeight={'bold'}>{props.email}</Text>.
            The link to reset your password expires in 24 hours.
          </Heading>
         
        </Stack>
   
      </Stack>
    </Flex>
  );
}