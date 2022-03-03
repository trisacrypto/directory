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

export default function PasswordReset() {
  return (
    <Flex
      minH={'100vh'}
      minWidth={'100vw'}
      align={'center'}
      justify={'center'}
      fontFamily={colors.font}
      fontSize={'xl'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={8} mx={'auto'} maxW={'lg'} py={12} px={6} width={'100%'}>
        <Stack align={'left'}>
          <Heading fontSize={'xl'}>Enter your email address.</Heading>
        </Stack>
        <Box
          rounded={'lg'}
          bg={useColorModeValue('white', 'transparent')}
          position={'relative'}
          bottom={5}
         >
          <Stack spacing={4} >
            <FormControl id="email">
              <Input type="email" height={'64px'} placeholder="Email Address"/>
            </FormControl>
            <Stack spacing={8} >
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
                Submit
              </Button>
            </Stack>
          </Stack>
        </Box>
      </Stack>
    </Flex>
  );
}
