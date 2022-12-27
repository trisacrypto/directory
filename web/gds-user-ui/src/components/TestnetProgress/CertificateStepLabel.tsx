import { FC, useEffect, useState } from 'react';
import {
  Box,
  Icon,
  Text,
  Heading,
  Stack,
  // Tooltip,
  Flex,
  useColorModeValue,
  useDisclosure,
  Button,
  Tooltip
} from '@chakra-ui/react';
import { FaCheckCircle, FaDotCircle, FaRegCircle } from 'react-icons/fa';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep, setHasReachSubmitStep } from 'application/store/stepper.slice';
import { findStepKey } from 'utils/utils';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import { useFormContext } from 'react-hook-form';
import useCertificateStepper from 'hooks/useCertificateStepper';
import InvalidFormPrompt from './InvalidFormPrompt';
import { useAppDispatch } from 'application/store';
export enum LCOLOR {
  'COMPLETE' = '#34A853',
  'PROGRESS' = '#5469D4',
  'SAVE' = '#F29C36',
  'INCOMPLETE' = '#C1C9D2',
  'NEXT' = '#E9E0E0',
  'ERROR' = '#dc2f02'
}
export enum LSTATUS {
  'COMPLETE' = 'complete',
  'PROGRESS' = 'progress',
  'SAVE' = 'save',
  'INCOMPLETE' = 'incomplete',
  'NEXT' = 'next',
  'ERROR' = 'error'
}
interface StepLabelProps {}
type TStepLabel = {
  color: string; // color of the icon
  hasError?: boolean; // status of the step
  icon: any; // icon of the step
};

// enum STEP {
//   BASIC_DETAILS = 1,
//   LEGAL_PERSON = 2,
//   CONTACTS = 3,
//   TRISA_IMPLEMENTATION = 4,
//   TRIXO_QUESTIONNAIRE = 5,
//   REVIEW = 6
// }

const CertificateStepLabel: FC<StepLabelProps> = () => {
  const dispatch = useAppDispatch();
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const textColor = useColorModeValue('#3C4257', '#F7F8FC');
  const { jumpToStep } = useCertificateStepper();
  const { isOpen, onClose, onOpen } = useDisclosure();
  const formContext = useFormContext();
  const [selectedStep, setSelectedStep] = useState<number>(currentStep);
  const [initialFormValues, setInitialFormValues] = useState<Record<string, any>>();

  useEffect(() => {
    setInitialFormValues(formContext.getValues());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const isStepCompleted = (step: number) => {
    const stepStatus = steps[step - 1]?.status;
    return stepStatus === 'complete' || stepStatus === 'progress';
  };

  // this function need some clean up
  const getLabel = (step: number): TStepLabel | undefined => {
    const s = findStepKey(steps, step);
    if (s && s?.length === 1) {
      if (s[0]?.status === LSTATUS.COMPLETE) {
        return {
          color: LCOLOR.COMPLETE,
          icon: FaCheckCircle
        };
      }
      if (s[0]?.status === LSTATUS.NEXT) {
        return {
          color: LCOLOR.NEXT,
          icon: FaCheckCircle
        };
      }

      if (s[0]?.status === LSTATUS.PROGRESS) {
        return {
          color: LCOLOR.PROGRESS,
          icon: FaDotCircle
        };
      }
      if (s[0]?.status === LSTATUS.SAVE) {
        return {
          color: LCOLOR.SAVE,
          icon: FaCheckCircle
        };
      }
      if (s[0]?.status === LSTATUS.ERROR) {
        return {
          color: LCOLOR.ERROR,
          icon: FaDotCircle,
          hasError: true
        };
      }
      if (s[0]?.status === LSTATUS.INCOMPLETE) {
        return {
          color: LCOLOR.INCOMPLETE,
          icon: FaRegCircle
        };
      }
    } else {
      return step === 1
        ? {
            color: LCOLOR.PROGRESS,
            icon: FaCheckCircle
          }
        : {
            color: LCOLOR.INCOMPLETE,
            icon: FaRegCircle
          };
    }
  };
  const isActiveStep = (step: number) => step === currentStep;

  const handleStepClick = (step: number) => () => {
    setSelectedStep(step);
    if (formContext.formState.isDirty) {
      onOpen();
    } else {
      dispatch(setHasReachSubmitStep({ hasReachSubmitStep: false }));
      jumpToStep(step);
    }
  };

  const handleContinueClick = () => {
    formContext.reset(initialFormValues);
    jumpToStep(selectedStep);
    onClose();
  };

  const stepLabels = [
    {
      label: t`Basic Details`
    },
    {
      label: t`Legal Person`
    },
    {
      label: t`Contacts`
    },
    {
      label: t`TRISA implementation`
    },
    {
      label: t`TRIXO Questionnaire`
    },
    {
      label: t`Review`
    }
  ];

  return (
    <>
      <Stack
        boxShadow="0px 10px 15px -3px rgba(0,0,0,0.1)"
        borderColor={'#C1C9D2'}
        borderRadius={8}
        borderWidth={1}
        bg={'white'}
        p={5}
        fontFamily={'Open Sans'}
        width="100%">
        <Box display={'flex'} justifyContent={'space-between'}>
          <Heading fontSize={['md', '2xl']} textTransform="capitalize">
            <Trans id="Certificate Progress">Certificate Progress</Trans>{' '}
          </Heading>
        </Box>
        <Flex gap={2}>
          {stepLabels.map((stepLabel, idx: number) => {
            const stepIndex = idx + 1;
            return (
              <Tooltip key={idx} label={stepLabel.label} gutter={0} hasArrow>
                <Button
                  bg="transparent"
                  display="block"
                  p={0}
                  width="100%"
                  height="100%"
                  _hover={{ bg: 'transparent' }}
                  disabled={!(() => isStepCompleted(stepIndex))()}
                  _disabled={{ opacity: 0.9, cursor: 'not-allowed' }}
                  onClick={handleStepClick(stepIndex)}>
                  <Stack spacing={1} width="100%">
                    <Box
                      h="1"
                      bg={getLabel(stepIndex)?.color}
                      borderRadius={'50px'}
                      width={'100%'}
                    />
                    <Stack
                      direction={{ base: 'column', md: 'row' }}
                      alignItems={{ base: 'center', lg: 'baseline' }}
                      spacing={{ base: 0, md: 1 }}>
                      <Box>
                        <Icon
                          as={getLabel(stepIndex)?.icon}
                          sx={{
                            path: {
                              fill: getLabel(stepIndex)?.color
                            },
                            verticalAlign: 'middle'
                          }}
                          verticalAlign={{ base: 'baseline', lg: 'middle' }}
                        />
                      </Box>
                      <Text
                        color={textColor}
                        fontSize={{ base: 'xs', md: 'sm' }}
                        fontWeight={isActiveStep(stepIndex) ? 'bold' : 'normal'}
                        textAlign="center"
                        noOfLines={1}>
                        {stepIndex} {stepLabel.label}
                      </Text>
                    </Stack>
                  </Stack>
                </Button>
              </Tooltip>
            );
          })}
          {/* <Button
            bg="transparent"
            display="block"
            p={0}
            width="100%"
            onClick={handleStepClick(STEP.BASIC_DETAILS)}
            disabled={!(() => isStepCompleted(STEP.BASIC_DETAILS))()}
            _disabled={{ opacity: 0.9, cursor: 'not-allowed' }}
            _hover={{ bg: 'transparent' }}>
            <Tooltip
              label={getLabel(STEP.BASIC_DETAILS)?.hasError && t`Missing required element`}
              placement="top"
              bg={'red'}>
              <Stack spacing={1} width="100%">
                <Box
                  h="1"
                  borderRadius={'50px'}
                  bg={getLabel(STEP.BASIC_DETAILS)?.color}
                  width={'100%'}
                  key={1}
                />
                <Stack
                  direction={{ base: 'column', md: 'row' }}
                  alignItems={['center']}
                  spacing={{ base: 0, md: 1 }}>
                  <Box>
                    <Icon
                      as={getLabel(STEP.BASIC_DETAILS)?.icon}
                      sx={{
                        path: {
                          fill: getLabel(STEP.BASIC_DETAILS)?.color
                        }
                      }}
                    />
                  </Box>
                  <Text
                    color={textColor}
                    fontWeight={isActiveStep(STEP.BASIC_DETAILS) ? 'bold' : 'normal'}
                    fontSize={'sm'}
                    textAlign="center">
                    1 <Trans id="Basic Details">Basic Details</Trans>
                  </Text>
                </Stack>
              </Stack>
            </Tooltip>
          </Button>

          <Button
            bg="transparent"
            display="block"
            p={0}
            width="100%"
            _hover={{ bg: 'transparent' }}
            disabled={!(() => isStepCompleted(STEP.LEGAL_PERSON))()}
            _disabled={{ opacity: 0.9, cursor: 'not-allowed' }}
            onClick={handleStepClick(STEP.LEGAL_PERSON)}>
            <Stack spacing={1} width="100%">
              <Box
                h="1"
                bg={getLabel(STEP.LEGAL_PERSON)?.color}
                borderRadius={'50px'}
                width={'100%'}
              />
              <Stack
                direction={{ base: 'column', md: 'row' }}
                alignItems={'center'}
                spacing={{ base: 0, md: 1 }}>
                <Box>
                  <Icon
                    as={getLabel(STEP.LEGAL_PERSON)?.icon}
                    sx={{
                      path: {
                        fill: getLabel(STEP.LEGAL_PERSON)?.color
                      },
                      verticalAlign: 'middle'
                    }}
                    verticalAlign={{ base: 'baseline', lg: 'middle' }}
                  />
                </Box>
                <Text
                  color={textColor}
                  fontSize={'sm'}
                  fontWeight={isActiveStep(STEP.LEGAL_PERSON) ? 'bold' : 'normal'}
                  textAlign="center">
                  2 <Trans id="Legal Person">Legal Person</Trans>
                </Text>
              </Stack>
            </Stack>
          </Button>

          <Button
            bg="transparent"
            display="block"
            p={0}
            width="100%"
            _hover={{ bg: 'transparent' }}
            disabled={!(() => isStepCompleted(STEP.CONTACTS))()}
            _disabled={{ opacity: 0.9, cursor: 'not-allowed' }}
            onClick={handleStepClick(STEP.CONTACTS)}>
            <Stack spacing={1} width="100%">
              <Box h="1" bg={getLabel(STEP.CONTACTS)?.color} width={'100%'} borderRadius={'50px'} />
              <Stack
                direction={{ base: 'column', md: 'row' }}
                alignItems={['center']}
                spacing={{ base: 0, md: 1 }}>
                <Box>
                  <Icon
                    as={getLabel(STEP.CONTACTS)?.icon}
                    sx={{
                      path: {
                        fill: getLabel(STEP.CONTACTS)?.color
                      }
                    }}
                  />
                </Box>
                <Text
                  color={textColor}
                  fontSize={'sm'}
                  fontWeight={isActiveStep(STEP.CONTACTS) ? 'bold' : 'normal'}
                  textAlign="center">
                  3 <Trans id="Contacts">Contacts</Trans>
                </Text>
              </Stack>
            </Stack>
          </Button>

          <Button
            bg="transparent"
            display="block"
            p={0}
            width="100%"
            _hover={{ bg: 'transparent' }}
            disabled={!(() => isStepCompleted(4))()}
            _disabled={{ opacity: 0.9, cursor: 'not-allowed' }}
            onClick={handleStepClick(4)}>
            <Stack spacing={1} width="100%">
              <Box h="1" bg={getLabel(4)?.color} width={'100%'} borderRadius={'50px'} />
              <Stack
                direction={{ base: 'column', md: 'row' }}
                alignItems={['center']}
                spacing={{ base: 0, md: 1 }}>
                <Box>
                  <Icon
                    as={getLabel(4)?.icon}
                    sx={{
                      path: {
                        fill: getLabel(4)?.color
                      }
                    }}
                  />
                </Box>
                <Text
                  color={textColor}
                  fontSize={'sm'}
                  fontWeight={isActiveStep(4) ? 'bold' : 'normal'}
                  textAlign="center">
                  4 <Trans id="TRISA implementation">TRISA implementation</Trans>
                </Text>
              </Stack>
            </Stack>
          </Button>

          <Button
            bg="transparent"
            display="block"
            p={0}
            width="100%"
            _hover={{ bg: 'transparent' }}
            disabled={!(() => isStepCompleted(STEP.TRIXO_QUESTIONNAIRE))()}
            _disabled={{ opacity: 0.9, cursor: 'not-allowed' }}
            onClick={handleStepClick(STEP.TRIXO_QUESTIONNAIRE)}>
            <Stack spacing={1} width="100%">
              <Box
                h="1"
                bg={getLabel(STEP.TRIXO_QUESTIONNAIRE)?.color}
                width={'100%'}
                borderRadius={'50px'}
              />
              <Stack
                direction={{ base: 'column', md: 'row' }}
                alignItems={['center']}
                spacing={{ base: 0, md: 1 }}>
                <Box>
                  <Icon
                    as={getLabel(STEP.TRIXO_QUESTIONNAIRE)?.icon}
                    sx={{
                      path: {
                        fill: getLabel(STEP.TRIXO_QUESTIONNAIRE)?.color
                      }
                    }}
                  />
                </Box>
                <Text
                  color={textColor}
                  fontSize={'sm'}
                  fontWeight={isActiveStep(STEP.TRIXO_QUESTIONNAIRE) ? 'bold' : 'normal'}
                  textAlign="center">
                  5 <Trans id="TRIXO Questionnaire">TRIXO Questionnaire</Trans>
                </Text>
              </Stack>
            </Stack>
          </Button>

          <Button
            bg="transparent"
            display="block"
            p={0}
            width="100%"
            _hover={{ bg: 'transparent' }}
            disabled={!(() => isStepCompleted(STEP.REVIEW))()}
            _disabled={{ opacity: 0.9, cursor: 'not-allowed' }}
            onClick={handleStepClick(STEP.REVIEW)}>
            <Stack spacing={1} width="100%">
              <Box h="1" bg={getLabel(STEP.REVIEW)?.color} width={'100%'} borderRadius={'50px'} />
              <Stack
                direction={{ base: 'column', md: 'row' }}
                alignItems={['center']}
                spacing={{ base: 0, md: 1 }}>
                <Box>
                  <Icon
                    as={getLabel(STEP.REVIEW)?.icon}
                    sx={{
                      path: {
                        fill: getLabel(STEP.REVIEW)?.color
                      }
                    }}
                  />
                </Box>
                <Text
                  color={textColor}
                  fontSize={'sm'}
                  fontWeight={isActiveStep(STEP.REVIEW) ? 'bold' : 'normal'}
                  textAlign="center">
                  6 <Trans id="Review">Review</Trans>
                </Text>
              </Stack>
            </Stack>
          </Button> */}
        </Flex>
        <InvalidFormPrompt
          isOpen={isOpen}
          onClose={onClose}
          handleContinueClick={handleContinueClick}
        />
      </Stack>
    </>
  );
};

export default CertificateStepLabel;
