import React, { useEffect, useState } from 'react';
import { Heading, Stack, Spinner } from '@chakra-ui/react';
import UserEmailVerification from 'components/Section/UserEmailVerification';
import UserEmailConfirmation from 'components/Section/UserEmailConfirmation';
import LandingLayout from 'layouts/LandingLayout';
import useQuery from 'hooks/useQuery';
import verifyService from './verify.service';
import AlertMessage from 'components/ui/AlertMessage';
import useAuth from 'hooks/useAuth';
import { useNavigate } from 'react-router-dom';
import Loader from 'components/Loader';
import TransparentLoader from 'components/Loader/TransparentLoader';
const VerifyPage: React.FC = () => {
  const query = useQuery();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<any>();
  const [result, setResult] = useState<any>(null);
  const [isRedirected, setIsRedirected] = useState(false);
  const vaspID = query.get('vaspID');
  const token = query.get('token');
  const registered_directory = query.get('registered_directory');
  const { isLoggedIn } = useAuth();
  const navigate = useNavigate();
  // to-do : should be improve later

  useEffect(() => {
    (async () => {
      try {
        if (vaspID && token && registered_directory) {
          const params = { vaspID, token, registered_directory };
          const response = await verifyService(params);
          if (!response.error) {
            setResult(response);
            // redirect to dashboard if user is logged in
            if (isLoggedIn()) {
              setTimeout(() => {
                setIsRedirected(true);
                navigate('/dashboard');
              }, 2000);
            }
          } else {
            console.error('Something went wrong');
            // setError(false)
          }
        } else {
          setError('Invalid params');
        }
      } catch (e: any) {
        if (!e.response?.data?.success) {
          setError(e.response?.data?.error);
        } else {
          // log error
          console.error('sorry something went wrong , please try again');
        }
      } finally {
        setIsLoading(false);
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [vaspID, token, registered_directory]);

  return (
    <LandingLayout>
      {isLoading && <Spinner size={'2xl'} />}
      {isRedirected && <TransparentLoader title={'Redirection to dashboard '} />}
      {result && (
        <AlertMessage message={result.message} status="success" title="Contact Verified " />
      )}
      {error && <AlertMessage message={error} status="error" />}
    </LandingLayout>
  );
};

export default VerifyPage;
