/* eslint-disable prefer-reflect */

import React, { useState, useEffect } from 'react';
import { Box, Heading, HStack, Icon, Stack, Text, useToast } from '@chakra-ui/react';

import BasicDetailsReview from './BasicDetailsReview';
import LegalPersonReview from './LegalPersonReview';
import TrisaImplementationReview from './TrisaImplementationReview';
import ContactsReview from './ContactsReview';
import TrixoReview from './TrixoReview';
import FormLayout from 'layouts/FormLayout';
import { RootStateOrAny, useSelector } from 'react-redux';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { hasStepError, mapTrixoFormForBff } from 'utils/utils';
import ReviewSubmit from 'components/ReviewSubmit';
import { registrationRequest } from 'modules/dashboard/certificate/service';
import { loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
const CertificateReview = () => {
  const { nextStep, previousStep } = useCertificateStepper();
  const toast = useToast();
  const steps: number = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const [isTestNetSent, setIsTestNetSent] = useState(false);
  const [isMainNetSent, setIsMainNetSent] = useState(false);
  const [result, setResult] = useState('');
  const handleSubmitRegister = async (event: React.FormEvent, network: string) => {
    event.preventDefault();
    try {
      const formValue = loadDefaultValueFromLocalStorage();
      const getMainnetObj = formValue.trisa_endpoint_mainnet;
      const getTestnetObj = formValue.trisa_endpoint_testnet;

      delete formValue.trisa_endpoint_mainnet;
      delete formValue.trisa_endpoint_testnet;
      if (network === 'testnet') {
        formValue.trisa_endpoint = getTestnetObj.endpoint;
        formValue.common_name = getTestnetObj.common_name;
      }
      if (network === 'mainnet') {
        formValue.trisa_endpoint = getMainnetObj.endpoint;
        formValue.common_name = getMainnetObj.common_name;
      }
      // clean form value before send to server
      // const formValueForBff = mapTrixoFormForBff(formValue);
      // console.log('formValueForBff', formValueForBff);
      const response = await registrationRequest(network, mapTrixoFormForBff(formValue));

      if (response.id || response.status === ' "SUBMITTED"') {
        if (network === 'testnet') {
          setIsTestNetSent(true);
          localStorage.setItem('isTestNetSent', 'true');
        }
        if (network === 'mainnet') {
          setIsMainNetSent(true);
          localStorage.setItem('isMainNetSent', 'true');
        }
        setResult(response);

        // toast({
        //   position: 'top-right',
        //   title: 'Success',
        //   description: response.message,
        //   status: 'success',
        //   duration: 5000,
        //   isClosable: true
        // });
      }
    } catch (err: any) {
      console.log('err', err?.response?.data);

      if (!err.response.data.success) {
        toast({
          position: 'top-right',
          title: 'Error Submitting Certificate',
          description: err.response.data.error,
          status: 'error',
          duration: 5000,
          isClosable: true
        });
      } else {
        console.log('something went wrong');
      }
    }
  };

  return (
    <>
      {!hasReachSubmitStep ? (
        <Stack spacing={7}>
          <HStack pt={10}>
            <Heading size="md"> Review </Heading>
            <Box>{/* <Icon as={InfoIcon} color="#F29C36" w={7} h={7} /> (not saved) */}</Box>
          </HStack>
          <FormLayout>
            <Text>
              Please review the information provided, edit as needed, and submit to complete the
              registration form. After the information is reviewed, you will be contacted to verify
              details. Once verified, your TestNet certificate will be issued.
            </Text>
          </FormLayout>
          <BasicDetailsReview />
          <LegalPersonReview />
          <ContactsReview />
          <TrisaImplementationReview />
          <TrixoReview />
        </Stack>
      ) : (
        <ReviewSubmit
          onSubmitHandler={handleSubmitRegister}
          isTestNetSent={isTestNetSent}
          isMainNetSent={isMainNetSent}
          result={result}
        />
      )}
    </>
  );
};

export default CertificateReview;
