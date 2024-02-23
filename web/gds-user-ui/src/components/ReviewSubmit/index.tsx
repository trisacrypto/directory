import { useEffect } from 'react';
import {
  Button,
  Heading,
  VStack,
  Stack,
  Text,
  useDisclosure,
  Box,
  Flex,
  Link,
  SimpleGrid,
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
import ConfirmationModal from 'components/ReviewSubmit/ConfirmationModal';
import { t, Trans } from '@lingui/macro';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { STEPPER_NETWORK } from 'utils/constants';

import WarningBox from 'components/WarningBox';
import { setHasReachSubmitStep } from 'application/store/stepper.slice';
import { useAppDispatch } from 'application/store';
import { StepsIndexes } from 'constants/steps';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import { StepEnum } from 'types/enums';
import useSubmissionStatus from 'modules/dashboard/registration/hooks/useSubmissionStatus';

interface ReviewSubmitProps {
  onSubmitHandler: (e: React.FormEvent, network: string) => void;
  isTestNetSent?: boolean;
  isMainNetSent?: boolean;
  result?: any;
  isTestNetSubmitting?: boolean;
  isMainNetSubmitting?: boolean;
  handleJumpToLastStep?: (e: React.FormEvent) => void;
}

const ReviewSubmit: React.FC<ReviewSubmitProps> = ({
  onSubmitHandler,
  isTestNetSent,
  isMainNetSent,
  result,
  isTestNetSubmitting,
  isMainNetSubmitting,
  handleJumpToLastStep,
}) => {
  const { certificateStep } = useFetchCertificateStep({ key: StepEnum.ALL });
  const { isOpen, onOpen, onClose } = useDisclosure();
  const isSent = isTestNetSent || isMainNetSent;
  const { jumpToStep } = useCertificateStepper();
  const { status } = useSubmissionStatus();
  const dispatch = useAppDispatch();
  
  // Display a warning box if the user has reached the submit step but not completed the
  // TRISA implementation step for one of the networks.
  const mainnetCommonName = certificateStep?.form?.mainnet?.common_name;
  const mainnetEndpoint = certificateStep?.form?.mainnet?.endpoint;
  const testnetCommonName = certificateStep?.form?.testnet?.common_name;
  const testnetEndpoint = certificateStep?.form?.testnet?.endpoint;

  const isMainnetNetworkIncomplete = !mainnetCommonName || !mainnetEndpoint;
  const isTestnetNetworkIncomplete = !testnetCommonName || !testnetEndpoint;

  const isTestNetSubmitted = status?.data?.testnet_submitted;
  const isMainNetSubmitted = status?.data?.mainnet_submitted;

  useEffect(() => {
    if (isSent) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isSent]);

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
              <Trans>Registration Submission</Trans>
            </Heading>
          </Stack>

          <FormLayout>
            <Text>
              <Trans>
                You must submit your registration for TestNet and MainNet separately.
              </Trans>{' '}
              <Text as="span" fontWeight="bold">
                <Trans>Note:</Trans>
              </Text>{' '}
              <Trans>
                You will receive 
              </Trans>{' '}
              <Text as="span" fontStyle={"italic"}>
                <Trans>two separate emails with confirmation links for each registration.</Trans>
              </Text>{' '}
              <Trans>
                You must click on each confirmation link to complete the registration process.
              </Trans>{' '}
              <Text as="span" fontWeight="bold">
                <Trans>
                Failure to click either confirmation will result in an incomplete registration.
                </Trans>
              </Text>
            </Text>
          </FormLayout>
          <SimpleGrid
            columns={{ base: 1, lg: 2 }}
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
                  <Trans>Click below to submit your</Trans>{' '}
                  <Text as="span" fontWeight={'bold'}>
                    <Trans>TestNet</Trans>
                  </Text>{' '}
                  <Trans>
                    registration. Upon submission, you will receive an email with a confirmation
                    link. You must click the confirmation link to complete the registration process.
                    Failure to click the confirmation link will result in an incomplete registration
                  </Trans>
                  .
                </Text>
                <Text>
                  <Trans>
                    A physical verification check in the form of a phone call
                  </Trans>{' '}
                  <Text as="span" fontWeight={'bold'}>
                    <Trans>is not required</Trans>
                  </Text>{' '}
                  <Trans>
                    for TestNet registration so your TestNet certificate will be issued upon review
                    by the validation team
                  </Trans>
                  .
                </Text>
                <Text>
                  <Trans>
                    If you would like to edit your registration form before submitting, please
                    return to the
                  </Trans>{' '}
                  <Link color="link" onClick={handleJumpToLastStep} fontWeight={'bold'}>
                    <Trans>Review section</Trans>
                  </Link>
                  .
                </Text>

                {isTestnetNetworkIncomplete ? (
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
                      <Trans>
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
                height="100%"
                mx="auto"
                pt="2"
                px="6"
                pb="6"
                alignItems={'center'}>
                <Button
                  bgColor="#ff7a59f0"
                  color="#fff"
                  size="lg"
                  whiteSpace="normal"
                  mt="auto"
                  py={{ base: '1rem', lg: '1.75rem' }}
                  width="100%"
                  boxShadow="lg"
                  _hover={{
                    bgColor: '#f55c35'
                  }}
                  isLoading={isTestNetSubmitting}
                  isDisabled={isTestnetNetworkIncomplete || isTestNetSubmitted }
                  data-testid="testnet-submit-btn"
                  onClick={(e) => {
                    onSubmitHandler(e, STEPPER_NETWORK.TESTNET);
                  }}>
                  {t`Submit TestNet Registration`}
                </Button>
              </Stack>
            </Stack>
            <Stack bg={'white'}>
              <Stack px={6} mb={5} pt={4}>
                <Heading size="sm" mt={2}>
                  <Trans>MAINNET SUBMISSION</Trans>
                </Heading>
                <Text>
                  <Trans>Click below to submit your</Trans>{' '}
                  <Text as="span" fontWeight={'bold'}>
                    <Trans>MainNet</Trans>
                  </Text>{' '}
                  <Trans>
                    registration. Upon submission, you will receive an email with a confirmation
                    link. You must click the confirmation link to complete the registration process.
                    Failure to click the confirmation link will result in an incomplete registration
                  </Trans>
                  .
                </Text>
                <Text>
                  <Trans>
                    A physical verification check in the form of a phone call
                  </Trans>{' '}
                  <Text as="span" fontWeight={'bold'}>
                    <Trans> is required</Trans>
                  </Text>{' '}
                  <Trans>
                    for MainNet registration so your MainNet certificate will be issued after the verification
                    phone call has been completed by the validation team.
                  </Trans>
                </Text>
                <Text>
                  <Trans>
                    If you would like to edit your registration form before submitting, please
                    return to the
                  </Trans>{' '}
                  <Link onClick={handleJumpToLastStep} color="link" fontWeight="bold">
                    <Trans>Review section</Trans>
                  </Link>
                  .
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
                      <Trans>
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
                height="100%"
                mx="auto"
                pt="2"
                alignItems={'center'}
                pb="6"
                px="6">
                <Button
                  bgColor="#23a7e0e8"
                  color="#fff"
                  size="lg"
                  mt="auto"
                  py={{ base: '1rem', lg: '1.75rem' }}
                  width="100%"
                  _hover={{
                    bgColor: '#189fda'
                  }}
                  isLoading={isMainNetSubmitting}
                  isDisabled={isMainnetNetworkIncomplete || isMainNetSubmitted }
                  whiteSpace="normal"
                  boxShadow="lg"
                  data-testid="mainnet-submit-btn"
                  onClick={(e) => {
                    onSubmitHandler(e, STEPPER_NETWORK.MAINNET);
                  }}
                  >
                  {t`Submit MainNet Registration`}
                </Button>
              </Stack>
            </Stack>
          </SimpleGrid>

          <Box alignSelf={'flex-start'} textAlign="center" mx={'auto'}>
            <Button
              bgColor="#fff"
              color="#1026F0"
              size="lg"
              py="2rem"
              whiteSpace="normal"
              boxShadow="lg"
              maxW="285px"
              width="100%"
              _hover={{
                bgColor: '#E6E6E6'
              }}
              data-cy="back-to-review-section"
              onClick={handleJumpToLastStep}>
              {t`Back to Review section`}
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
