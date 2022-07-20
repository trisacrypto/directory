import React, { useState } from 'react';
import { useToast } from '@chakra-ui/react';
import PasswordReset from 'components/Section/PasswordReset';
import * as Sentry from '@sentry/browser';
import LandingLayout from 'layouts/LandingLayout';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import { t } from '@lingui/macro';

const ResetPassword: React.FC = () => {
  const [isLoading, setIsloading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>('');
  const { auth0ResetPassword } = useCustomAuth0();
  const toast = useToast();
  const handleResetPassword = async (data: any) => {
    setIsloading(true);
    try {
      const option = {
        email: data.username,
        connection: 'Username-Password-Authentication'
      };

      const response: any = await auth0ResetPassword(option);
      setIsloading(false);
      if (response) {
        const content = t`Thank you. We have sent instructions to reset your password to ${data.username}. The link to reset your password expires in 24 hours.`;
        setMessage(content);
      }
    } catch (err: any) {
      setIsloading(false);
      Sentry.captureException(err);
    }
  };

  return (
    <LandingLayout>
      <PasswordReset isLoading={isLoading} handleSubmit={handleResetPassword} message={message} />
    </LandingLayout>
  );
};

export default ResetPassword;
