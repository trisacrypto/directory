/* eslint-disable prefer-reflect */
// TO DO: refactor certificate stepper to use react-query to fetch data and handle loading state
// TO DO: Write clean code for this component and make it more easily testable

import React, { lazy, Suspense, useState, useEffect } from 'react';
import { useToast } from '@chakra-ui/react';

import { RootStateOrAny, useSelector, useDispatch } from 'react-redux';
import ReviewSubmit from 'components/ReviewSubmit';
import { t } from '@lingui/macro';
import {
  submitMainnetRegistration,
  submitTestnetRegistration
} from 'modules/dashboard/registration/service';
import useCertificateStepper from 'hooks/useCertificateStepper';
import Loader from 'components/Loader';
import { getRefreshToken } from 'utils/auth0.helper';
import { STEPPER_NETWORK } from 'utils/constants';
import { setUserName } from 'modules/auth/login/user.slice';
const ReviewsSummary = lazy(() => import('./ReviewsSummary'));

const CertificateReview = () => {
  const toast = useToast();
  const dispatch = useDispatch();
  const { testnetSubmissionState, mainnetSubmissionState } = useCertificateStepper();

  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const [isTestNetSent, setIsTestNetSent] = useState(false);
  const [isMainNetSent, setIsMainNetSent] = useState(false);
  const [isTestNetSubmitting, setIsTestNetSubmitting] = useState(false);
  const [isMainNetSubmitting, setIsMainNetSubmitting] = useState(false);
  const [result, setResult] = useState<any>('');
  // refactor this to use react-query
  const handleSubmitRegister = async (event: React.FormEvent, network: string) => {
    event.preventDefault();
    try {
      if (network === STEPPER_NETWORK.TESTNET) {
        setIsTestNetSubmitting(true);
        const response = await submitTestnetRegistration();
        if (response.status === 200) {
          await getRefreshToken(response?.data?.refresh_token);
          setIsTestNetSubmitting(false);
          setIsTestNetSent(true);
          testnetSubmissionState();

          setResult(response?.data);
        }
      }
      if (network === STEPPER_NETWORK.MAINNET) {
        setIsMainNetSubmitting(true);
        const response = await submitMainnetRegistration();
        if (response?.status === 200) {
          await getRefreshToken(response?.data?.refresh_token);
          setIsMainNetSubmitting(false);

          setIsMainNetSent(true);
          mainnetSubmissionState();
          setResult(response?.data);
        }
      }
    } catch (err: any) {
      setIsMainNetSubmitting(false);
      setIsTestNetSubmitting(false);

      if (!err?.response?.data?.success) {
        toast({
          position: 'top-right',
          title: t`Error Submitting Certificate`,
          description: t`${err?.response?.data?.error}`,
          status: 'error',
          duration: 5000,
          isClosable: true
        });
      } else {
        console.error('something went wrong');
      }
    }
  };
  useEffect(() => {
    if (isTestNetSent || isMainNetSent) {
      dispatch(setUserName(result?.organization_name));
    }
  }, [isTestNetSent, isMainNetSent, result, dispatch]);

  if (!hasReachSubmitStep) {
    return <ReviewsSummary />;
  }

  return (
    <Suspense fallback={<Loader />}>
      <ReviewSubmit
        onSubmitHandler={handleSubmitRegister}
        isTestNetSent={isTestNetSent}
        isMainNetSent={isMainNetSent}
        result={result}
        isTestNetSubmitting={isTestNetSubmitting}
        isMainNetSubmitting={isMainNetSubmitting}
      />
    </Suspense>
  );
};

export default CertificateReview;
