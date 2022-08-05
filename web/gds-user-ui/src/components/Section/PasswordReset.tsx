import { Flex, Stack, useColorModeValue, Text } from '@chakra-ui/react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useEffect, useState } from 'react';
import SuccessMessage from 'components/ui/SuccessMessage';
import { colors } from '../../utils/theme';
import { Trans } from '@lingui/react';
import PasswordResetForm from './PasswordResetForm';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import * as Sentry from '@sentry/browser';
import { t } from '@lingui/macro';
import * as Yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup';

export type ResetPasswordFormValues = {
  email: string;
};

const validationSchema = Yup.object().shape({
  email: Yup.string().email('Email is invalid').required('Email is required')
});

const PasswordReset = () => {
  const methods = useForm<ResetPasswordFormValues>({
    defaultValues: {
      email: ''
    },
    resolver: yupResolver(validationSchema)
  });
  const { resetField } = methods;
  const [message, setMessage] = useState<string>('');
  const { auth0ResetPassword } = useCustomAuth0();

  const handleResetPassword: SubmitHandler<ResetPasswordFormValues> = async (
    data: ResetPasswordFormValues
  ) => {
    try {
      const option = {
        email: data.email,
        connection: 'Username-Password-Authentication'
      };

      const response: any = await auth0ResetPassword(option);
      if (response) {
        const content = t`Thank you. We have sent instructions to reset your password to ${data.email}. The link to reset your password expires in 24 hours.`;
        setMessage(content);
      }
    } catch (err: any) {
      Sentry.captureException(err);
    }
  };

  useEffect(() => {
    if (message) {
      resetField('email');
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [message]);

  return (
    <FormProvider {...methods}>
      <Flex
        align={'center'}
        justify={'center'}
        fontFamily={colors.font}
        color={useColorModeValue('gray.600', 'white')}
        fontSize={'xl'}
        bg={useColorModeValue('white', 'gray.800')}>
        <Stack spacing={8} mx={'auto'} maxW={'lg'} py={12} px={6} width={'100%'}>
          {message && <SuccessMessage message={message} handleClose={() => {}} />}
          <Stack align={'left'}>
            <Text fontSize="lg" mb={3} fontWeight="bold">
              Follow the instructions below to reset your TRISA password
            </Text>
            <Text fontSize={'sm'}>
              <Trans id="Enter your email address">Enter your email address</Trans>
            </Text>
          </Stack>

          <PasswordResetForm onSubmit={handleResetPassword} />
        </Stack>
      </Flex>
    </FormProvider>
  );
};

export default PasswordReset;
