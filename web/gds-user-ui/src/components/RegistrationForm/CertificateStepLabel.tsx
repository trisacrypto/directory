import { FC, useState } from 'react';
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
import { TStep, setHasReachSubmitStep, setIsDirty } from 'application/store/stepper.slice';
import { findStepKey } from 'utils/utils';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';

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
  const { jumpToStep, getIsDirtyState } = useCertificateStepper();
  const { isOpen, onClose, onOpen } = useDisclosure();

  const [selectedStep, setSelectedStep] = useState<number>(currentStep);

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
    dispatch(setHasReachSubmitStep({ hasReachSubmitStep: false }));
    setSelectedStep(step);
    // check if the current step is added to the stepper
    // if not then add it on steps state

    if (getIsDirtyState()) {
      onOpen();
    } else {
      jumpToStep(step);
    }
  };

  const handleContinueClick = () => {
    // set dirty current state to false
    dispatch(setIsDirty({ isDirty: false }));
    // updateStepStatusToIncomplete();

    jumpToStep(selectedStep);
    onClose();
  };

  const getNextStepBtn = () => {
    let content = '' as string;
    if (currentStep - selectedStep < 0) {
      content = `Save & Next`;
    } else {
      content = `Save & Previous`;
    }

    return content;
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
          {stepLabels.map((stepLabel: any, idx: number) => {
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
        </Flex>
        <InvalidFormPrompt
          isOpen={isOpen}
          onClose={onClose}
          handleContinueClick={handleContinueClick}
          nextStepBtnContent={getNextStepBtn()}
        />
      </Stack>
    </>
  );
};

export default CertificateStepLabel;
