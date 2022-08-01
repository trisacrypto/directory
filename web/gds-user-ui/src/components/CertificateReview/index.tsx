/* eslint-disable prefer-reflect */

import React, { useState } from 'react';
import { useToast } from '@chakra-ui/react';

import { RootStateOrAny, useSelector } from 'react-redux';
import ReviewSubmit from 'components/ReviewSubmit';
import { registrationRequest } from 'modules/dashboard/certificate/service';
import { loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
import { t } from '@lingui/macro';
import ReviewsSummary from './ReviewsSummary';
import { mapTrixoFormForBff } from 'utils/utils';
import {
  submitMainnetRegistration,
  submitTestnetRegistration
} from 'modules/dashboard/registration/service';
const CertificateReview = () => {
  const toast = useToast();

  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const [isTestNetSent, setIsTestNetSent] = useState(false);
  const [isMainNetSent, setIsMainNetSent] = useState(false);
  const [result, setResult] = useState('');
  const handleSubmitRegister = async (event: React.FormEvent, network: string) => {
    event.preventDefault();
    try {
      // const formValue = loadDefaultValueFromLocalStorage();
      // const getMainnetObj = formValue.trisa_endpoint_mainnet;
      // const getTestnetObj = formValue.trisa_endpoint_testnet;

      // delete formValue.trisa_endpoint_mainnet;
      // delete formValue.trisa_endpoint_testnet;
      // if (network === 'testnet') {
      //   formValue.trisa_endpoint = getTestnetObj.endpoint;
      //   formValue.common_name = getTestnetObj.common_name;
      // }
      // if (network === 'mainnet') {
      //   formValue.trisa_endpoint = getMainnetObj.endpoint;
      //   formValue.common_name = getMainnetObj.common_name;
      // }

      if (network === 'testnet') {
        const response = await submitTestnetRegistration();
        console.log('[response testnet]', response);
        if (response.status === 200) {
          setIsTestNetSent(true);
          setResult(response?.data);
        }
      }
      if (network === 'mainnet') {
        const response = await submitMainnetRegistration();
        console.log('[response mainnet]', response);
        if (response?.status === 200) {
          setIsMainNetSent(true);
          setResult(response?.data);
        }
      }
    } catch (err: any) {
      console.log('[err catched 0]', err);

      if (!err?.response?.data?.success) {
        console.log('[err catched]', err?.response);
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

  if (!hasReachSubmitStep) {
    return <ReviewsSummary />;
  }

  return (
    <ReviewSubmit
      onSubmitHandler={handleSubmitRegister}
      isTestNetSent={isTestNetSent}
      isMainNetSent={isMainNetSent}
      result={result}
    />
  );
};

export default CertificateReview;
