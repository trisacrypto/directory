import React, { useEffect, useState } from 'react';
import { Spinner } from '@chakra-ui/react';
import LandingLayout from 'layouts/LandingLayout';
import useQuery from 'hooks/useQuery';
import verifyService from './verify.service';
import AlertMessage from 'components/ui/AlertMessage';
import useAuth from 'hooks/useAuth';
import { useNavigate } from 'react-router-dom';
import TransparentLoader from 'components/Loader/TransparentLoader';
import { upperCaseFirstLetter } from 'utils/utils';
const VerifyPage: React.FC = () => {
  const { vaspID, token, registered_directory } = useQuery();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<any>();
  const [result, setResult] = useState<any>(null);
  const [isRedirected, setIsRedirected] = useState<boolean>(false);
  const { isLoggedIn } = useAuth();
  const navigate = useNavigate();

  // to-do : should be improve later -  Also improve the alert message to be more user friendly

  useEffect(() => {
    (async () => {
      try {
        if (vaspID && token && registered_directory) {
          setIsRedirected(false);
          const params = { vaspID, token, registered_directory };
          const response = await verifyService(params);
          if (!response.error) {
            setResult(response);
            // redirect to dashboard if user is logged in
            if (isLoggedIn) {
              setTimeout(() => {
                setIsRedirected(true);
                navigate('/dashboard');
              }, 1000);
            }
          } else {
            console.error('Something went wrong');
            // setError(false)
          }
        } else {
          setError('Invalid params');
        }
      } catch (e: any) {
        if (!e?.data?.success) {
          setError(upperCaseFirstLetter(e?.data?.error));
        } else {
          // log error
          console.error('Sorry something went wrong, please try again.');
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
      {isRedirected && <TransparentLoader title={'Redirecting to the Dashboard...'} />}
      {result && (
        <AlertMessage message={result.message} status="success" title="Contact Verified " />
      )}
      {error && <AlertMessage message={error} status="error" />}
    </LandingLayout>
  );
};

export default VerifyPage;
