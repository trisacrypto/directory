import { useState, useEffect } from 'react';
import {
  Button,
  Heading,
  VStack,
  Stack,
  Text,
  useDisclosure,
  Box,
  Flex,
  Link
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
import ConfirmationModal from 'components/ReviewSubmit/ConfirmationModal';
import { t, Trans } from '@lingui/macro';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { useSelector } from 'react-redux';
// import { useNavigate } from 'react-router-dom';
import { STEPPER_NETWORK } from 'utils/constants';
import {
  getTestNetSubmittedStatus,
  getMainNetSubmittedStatus
} from 'application/store/selectors/stepper';

import WarningBox from 'components/WarningBox';
import { setHasReachSubmitStep } from 'application/store/stepper.slice';
import { useAppDispatch } from 'application/store';
import { StepsIndexes } from 'constants/steps';

interface ReviewSubmitProps {
  onSubmitHandler: (e: React.FormEvent, network: string) => void;
  isTestNetSent?: boolean;
  isMainNetSent?: boolean;
  result?: any;
  isTestNetSubmitting?: boolean;
  isMainNetSubmitting?: boolean;
}

const ReviewSubmit: React.FC<ReviewSubmitProps> = ({
  onSubmitHandler,
  isTestNetSent,
  isMainNetSent,
  result,
  isTestNetSubmitting,
  isMainNetSubmitting
}) => {
  const isTestNetSubmitted: boolean = useSelector(getTestNetSubmittedStatus);
  const isMainNetSubmitted: boolean = useSelector(getMainNetSubmittedStatus);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const isSent = isTestNetSent || isMainNetSent;
  const [testnet, setTestnet] = useState(false);
  const [mainnet, setMainnet] = useState(false);
  const { jumpToLastStep, jumpToStep } = useCertificateStepper();
  // const navigate = useNavigate();
  const dispatch = useAppDispatch();

  const isTestnetNetworkFieldsIncomplete = false;
  const isMainnetNetworkIncomplete = false;
  useEffect(() => {
    if (isTestNetSubmitted) {
      setTestnet(true);
    }
    if (isMainNetSubmitted) {
      setMainnet(true);
    }
  }, [isTestNetSubmitted, isMainNetSubmitted]);
  useEffect(() => {
    if (isSent) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isTestNetSent, isMainNetSent]);

  const handleJumpToLastStep = () => {
    jumpToLastStep();
    // navigate('/dashboard/certificate/registration');
  };

  const handleJumpToTrisaImplementationStep = () => {
    dispatch(setHasReachSubmitStep({ hasReachSubmitStep: false }));
    jumpToStep(StepsIndexes.TRISA_IMPLEMENTATION);
  };

  return (
    <>
      <Flex>
        <VStack mt="2rem">
          <Stack align="start" w="full">
            <Heading size="md" pr={3} ml={2}>
              <Trans id="Registration Submission">Registration Submission</Trans>
            </Heading>
          </Stack>

          <FormLayout>
            <Text>
              <Trans id="You must submit your registration for TestNet and MainNet separately.">
                You must submit your registration for TestNet and MainNet separately.
              </Trans>{' '}
              <Text as="span" fontWeight="bold">
                <Trans id="Note:">Note:</Trans>
              </Text>{' '}
              <Trans id="You will receive two separate emails with confirmation links for each registration. You must click on each confirmation link to complete the registration process.">
                You will receive two separate emails with confirmation links for each registration.
                You must click on each confirmation link to complete the registration process.{' '}
              </Trans>
              <Text as="span" fontWeight="bold">
                <Trans id=" Failure to click either confirmation will result in an incomplete registration">
                  Failure to click either confirmation will result in an incomplete registration.
                </Trans>
              </Text>
            </Text>
          </FormLayout>
          <Stack
            direction={['column', 'row']}
            justifyContent="space-around"
            py={14}
            width="100%"
            spacing={10}>
            <Stack bg={'white'}>
              <Stack px={6} mb={5} pt={4}>
                <Heading size="sm" mt={2}>
                  <Trans id="TESTNET SUBMISSION">TESTNET SUBMISSION</Trans>
                </Heading>
                <Text>
                  <Trans id="Click below to submit your">Click below to submit your</Trans>{' '}
                  <Text as="span" fontWeight={'bold'}>
                    <Trans id="TestNet">TestNet</Trans>
                  </Text>{' '}
                  <Trans id="registration. Upon submission, you will receive an email with a confirmation link. You must click the confirmation link to complete the registration process. Failure to click the confirmation link will result in an incomplete registration">
                    registration. Upon submission, you will receive an email with a confirmation
                    link. You must click the confirmation link to complete the registration process.
                    Failure to click the confirmation link will result in an incomplete registration
                  </Trans>
                  .
                </Text>
                <Text>
                  <Trans id="A physical verification check in the form of a phone call">
                    A physical verification check in the form of a phone call
                  </Trans>{' '}
                  <Text as="span" fontWeight={'bold'}>
                    <Trans id="is not required">is not required</Trans>
                  </Text>{' '}
                  <Trans id="for TestNet registration so your TestNet certificate will be issued upon review by the validation team">
                    for TestNet registration so your TestNet certificate will be issued upon review
                    by the validation team
                  </Trans>
                  .
                </Text>
                <Text>
                  <Trans id="If you would like to edit your registration form before submitting, please return to the">
                    If you would like to edit your registration form before submitting, please
                    return to the
                  </Trans>{' '}
                  <Link color="link" onClick={handleJumpToLastStep} fontWeight={'bold'}>
                    <Trans id="Review page">Review page</Trans>
                  </Link>
                  .
                </Text>

                {isTestnetNetworkFieldsIncomplete ? (
                  <WarningBox>
                    <Text>
                      <Trans>
                        If you would like to register for TestNet, please provide a{' '}
                        <Link
                          color="#1F4CED"
                          fontWeight={500}
                          onClick={handleJumpToTrisaImplementationStep}>
                          TestNet Endpoint and Common Name
                        </Link>
                        .
                      </Trans>
                    </Text>
                    <Text>
                      <Trans id="Please note that TestNet and MainNet are separate networks that require different X.509 Identity Certificates.">
                        Please note that TestNet and MainNet are separate networks that require
                        different X.509 Identity Certificates.
                      </Trans>
                    </Text>
                  </WarningBox>
                ) : null}
              </Stack>
              <Stack
                alignContent={'center'}
                justifyContent={'center'}
                mx="auto"
                pb={4}
                alignItems={'center'}>
                <Button
                  bgColor="#ff7a59f0"
                  color="#fff"
                  isDisabled={testnet || isTestnetNetworkFieldsIncomplete}
                  data-testid="testnet-submit-btn"
                  size="lg"
                  py="2.5rem"
                  isLoading={isTestNetSubmitting}
                  whiteSpace="normal"
                  maxW="200px"
                  width="100%"
                  boxShadow="lg"
                  onClick={(e) => {
                    onSubmitHandler(e, STEPPER_NETWORK.TESTNET);
                  }}
                  _hover={{
                    bgColor: '#f55c35'
                  }}>
                  {t`Submit TestNet Registration`}
                </Button>
              </Stack>
            </Stack>
            <Stack bg={'white'}>
              <Stack px={6} mb={5} pt={4}>
                <Heading size="sm" mt={2}>
                  <Trans id="MAINNET SUBMISSION">MAINNET SUBMISSION</Trans>
                </Heading>
                <Text>
                  <Trans id="Click below to submit your">Click below to submit your</Trans>{' '}
                  <Text as="span" fontWeight={'bold'}>
                    <Trans id="MainNet">MainNet</Trans>
                  </Text>{' '}
                  <Trans id="registration. Upon submission, you will receive an email with a confirmation link. You must click the confirmation link to complete the registration process. Failure to click the confirmation link will result in an incomplete registration">
                    registration. Upon submission, you will receive an email with a confirmation
                    link. You must click the confirmation link to complete the registration process.
                    Failure to click the confirmation link will result in an incomplete registration
                  </Trans>
                  .
                </Text>
                <Text>
                  <Trans id="A physical verification check in the form of a phone call">
                    physical verification check in the form of a phone call
                  </Trans>{' '}
                  <Text as="span" fontWeight={'bold'}>
                    <Trans id="is required"> is required</Trans>
                  </Text>{' '}
                  <Trans id=" for MainNet registration so your MainNet certificate will be issued upon review by the validation team">
                    for MainNet registration so your MainNet certificate will be issued upon review
                    by the validation team
                  </Trans>
                </Text>
                <Text>
                  <Trans id="If you would like to edit your registration form before submitting, please return to the">
                    If you would like to edit your registration form before submitting, please
                    return to the
                  </Trans>{' '}
                  <Link onClick={handleJumpToLastStep} color="link" fontWeight="bold">
                    <Trans id="Review page">Review page</Trans>
                  </Link>
                </Text>
                {isMainnetNetworkIncomplete ? (
                  <WarningBox>
                    <Text>
                      <Trans>
                        If you would like to register for MainNet, please provide a{' '}
                        <Link
                          color="#1F4CED"
                          fontWeight={500}
                          onClick={handleJumpToTrisaImplementationStep}>
                          MainNet Endpoint and Common Name
                        </Link>
                        .
                      </Trans>
                    </Text>
                    <Text>
                      <Trans id="Please note that TestNet and MainNet are separate networks that require different X.509 Identity Certificates.">
                        Please note that TestNet and MainNet are separate networks that require
                        different X.509 Identity Certificates.
                      </Trans>
                    </Text>
                  </WarningBox>
                ) : null}
              </Stack>
              <Stack
                alignContent={'center'}
                justifyContent={'center'}
                mx="auto"
                alignItems={'center'}
                pb={4}>
                <Button
                  bgColor="#23a7e0e8"
                  color="#fff"
                  size="lg"
                  py="2.5rem"
                  isLoading={isMainNetSubmitting}
                  isDisabled={mainnet || isMainnetNetworkIncomplete}
                  whiteSpace="normal"
                  boxShadow="lg"
                  data-testid="mainnet-submit-btn"
                  maxW="200px"
                  onClick={(e) => {
                    onSubmitHandler(e, STEPPER_NETWORK.MAINNET);
                  }}
                  width="100%"
                  _hover={{
                    bgColor: '#189fda'
                  }}>
                  {t`Submit MainNet Registration`}
                </Button>
              </Stack>
            </Stack>
          </Stack>

          <Box alignItems={'center'} textAlign="center" mx={'auto'}>
            <Button
              bgColor="#555151"
              color="#fff"
              onClick={handleJumpToLastStep}
              size="lg"
              py="2.5rem"
              whiteSpace="normal"
              boxShadow="lg"
              maxW="200px"
              width="100%"
              _hover={{
                bgColor: '#555151'
              }}>
              {t`Back to Review Page`}
            </Button>
          </Box>
        </VStack>
        {isSent && (
          <ConfirmationModal
            isOpen={isOpen}
            onClose={onClose}
            id={result?.id}
            pkcs12password={result?.pkcs12password}
            message={result?.message}
            status={result?.status}
            size={'xl'}
          />
        )}
      </Flex>
    </>
  );
};

export default ReviewSubmit;
