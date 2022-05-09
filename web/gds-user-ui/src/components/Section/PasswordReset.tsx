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
  Image
} from '@chakra-ui/react';
import { useForm } from 'react-hook-form';
import React, { useEffect, useState } from 'react';
import SuccessMessage from 'components/ui/SuccessMessage';
import { colors } from '../../utils/theme';

interface PasswordResetProps {
  handleSubmit: (data: any) => void;
  isLoading: boolean;
  isError?: any;
  message?: string;
}
const PasswordReset: React.FC<PasswordResetProps> = (props) => {
  const { register, handleSubmit, resetField } = useForm();

  useEffect(() => {
    if (props.message) {
      resetField('username');
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.message]);

  return (
    <Flex
      minWidth={'100vw'}
      align={'center'}
      justify={'center'}
      fontFamily={colors.font}
      fontSize={'xl'}
      marginTop={'10vh'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={8} mx={'auto'} maxW={'lg'} py={12} px={6} width={'100%'}>
        {props.message && <SuccessMessage message={props.message} handleClose={() => {}} />}
        <Stack align={'left'}>
          <Heading fontSize={'xl'}>Enter your email address.</Heading>
        </Stack>

        <Box
          rounded={'lg'}
          bg={useColorModeValue('white', 'transparent')}
          position={'relative'}
          bottom={5}>
          <form onSubmit={handleSubmit(props.handleSubmit)}>
            <Stack spacing={4}>
              <FormControl id="email">
                <Input
                  type="email"
                  height={'64px'}
                  {...register('username')}
                  placeholder="Email Address"
                />
              </FormControl>
              <Stack spacing={8}>
                <Button
                  bg={colors.system.blue}
                  color={'white'}
                  height={'57px'}
                  isLoading={props.isLoading}
                  type="submit"
                  w={['full', '50%']}
                  _hover={{
                    bg: '#10aaed'
                  }}
                  _focus={{
                    borderColor: 'transparent'
                  }}>
                  Submit
                </Button>
              </Stack>
            </Stack>
          </form>
        </Box>
      </Stack>
    </Flex>
  );
};

export default PasswordReset;
