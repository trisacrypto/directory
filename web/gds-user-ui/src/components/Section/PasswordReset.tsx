import {
  Flex,
  Box,
  FormControl,
  Input,
  Stack,
  Button,
  Heading,
  useColorModeValue,
  Image,
  Text
} from '@chakra-ui/react';
import { useForm } from 'react-hook-form';
import React, { useEffect } from 'react';
import SuccessMessage from 'components/ui/SuccessMessage';
import { colors } from '../../utils/theme';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';

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
      align={'center'}
      justify={'center'}
      fontFamily={colors.font}
      color={useColorModeValue('gray.600', 'white')}
      fontSize={'xl'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={8} mx={'auto'} maxW={'lg'} py={12} px={6} width={'100%'}>
        {props.message && <SuccessMessage message={props.message} handleClose={() => {}} />}
        <Stack align={'left'}>
          <Text fontSize="lg" mb={3} fontWeight="bold">
            Follow the instructions below to reset your TRISA password
          </Text>
          <Text fontSize={'sm'}>
            <Trans id="Enter your email address">Enter your email address</Trans>
          </Text>
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
                  size="lg"
                  {...register('username')}
                  placeholder={t`Email Address`}
                />
              </FormControl>
              <Button
                display="block"
                alignSelf="start"
                px={16}
                bg="blue"
                color={'white'}
                isLoading={props.isLoading}
                type="submit"
                _hover={{
                  bg: '#10aaed'
                }}
                _focus={{
                  borderColor: 'transparent'
                }}>
                <Trans id="Submit">Submit</Trans>
              </Button>
            </Stack>
          </form>
        </Box>
      </Stack>
    </Flex>
  );
};

export default PasswordReset;
