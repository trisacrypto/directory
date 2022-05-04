import React, { FC, useEffect } from 'react';
import { HStack, Box, Icon, Text, Heading, Stack, Grid, Button, Tooltip } from '@chakra-ui/react';
import { FaCheckCircle, FaDotCircle, FaRegCircle } from 'react-icons/fa';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { addStep, setCurrentStep, setStepStatus, TStep } from 'application/store/stepper.slice';
import { findStepKey } from 'utils/utils';
import { IconType } from 'react-icons/lib';
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
const CertificateStepLabel: FC<StepLabelProps> = (props) => {
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);

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

  return (
    <>
      <Box
        position={'relative'}
        bg={'white'}
        height={'96px'}
        pt={5}
        boxShadow="0 24px 50px rgba(55,65, 81, 0.25) "
        borderColor={'#C1C9D2'}
        borderRadius={8}
        borderWidth={1}
        // mt={10}
        // mx={5}
        px={5}
        fontFamily={'Open Sans'}>
        <Box pb={2} display={'flex'} justifyContent={'space-between'}>
          <Heading fontSize={20}> Certificate Progress </Heading>
        </Box>
        <Grid templateColumns="repeat(6, 1fr)" gap={2}>
          <Tooltip
            label={getLabel(1)?.hasError && 'Missing required element'}
            placement="top"
            bg={'red'}>
            <Box w="70px" h="1" borderRadius={50} bg={getLabel(1)?.color} width={'100%'} key={1}>
              <HStack>
                <Box pt={3}>
                  <Icon
                    as={getLabel(1)?.icon}
                    sx={{
                      path: {
                        fill: getLabel(1)?.color
                      }
                    }}
                  />
                </Box>
                <Text
                  pt={2}
                  color={'#3C4257'}
                  fontWeight={isActiveStep(1) ? 'bold' : 'normal'}
                  fontSize={'0.8em'}>
                  Basic Details
                </Text>
              </HStack>
            </Box>
          </Tooltip>
          <Box w="70px" h="1" bg={getLabel(2)?.color} width={'100%'}>
            <HStack>
              <Box pt={3}>
                <Icon
                  as={getLabel(2)?.icon}
                  sx={{
                    path: {
                      fill: getLabel(1)?.color
                    }
                  }}
                />
              </Box>
              <Text
                pt={2}
                color={'#3C4257'}
                fontSize={'0.8em'}
                fontWeight={isActiveStep(2) ? 'bold' : 'normal'}>
                Legal Person
              </Text>
            </HStack>
          </Box>

          <Box w="70px" h="1" bg={getLabel(3)?.color} width={'100%'}>
            <HStack>
              <Box pt={3}>
                <Icon
                  as={getLabel(3)?.icon}
                  sx={{
                    path: {
                      fill: getLabel(1)?.color
                    }
                  }}
                />
              </Box>
              <Text
                pt={2}
                color={'#3C4257'}
                fontSize={'0.8em'}
                fontWeight={isActiveStep(3) ? 'bold' : 'normal'}>
                Contacts
              </Text>
            </HStack>
          </Box>

          <Box w="70px" h="1" bg={getLabel(4)?.color} width={'100%'}>
            <HStack>
              <Box pt={3}>
                <Icon
                  as={getLabel(4)?.icon}
                  sx={{
                    path: {
                      fill: getLabel(1)?.color
                    }
                  }}
                />
              </Box>
              <Text
                pt={2}
                color={'#3C4257'}
                fontSize={'0.8em'}
                fontWeight={isActiveStep(4) ? 'bold' : 'normal'}>
                Trisa implementation
              </Text>
            </HStack>
          </Box>

          <Box w="70px" h="1" bg={getLabel(5)?.color} width={'100%'}>
            <HStack>
              <Box pt={3}>
                <Icon
                  as={getLabel(5)?.icon}
                  sx={{
                    path: {
                      fill: getLabel(1)?.color
                    }
                  }}
                />
              </Box>
              <Text
                pt={2}
                color={'#3C4257'}
                fontSize={'0.8em'}
                fontWeight={isActiveStep(5) ? 'bold' : 'normal'}>
                TRIXO Questionnaire
              </Text>
            </HStack>
          </Box>

          <Box w="70px" h="1" bg={getLabel(6)?.color} width={'100%'}>
            <HStack>
              <Box pt={3}>
                <Icon
                  as={getLabel(6)?.icon}
                  sx={{
                    path: {
                      fill: getLabel(1)?.color
                    }
                  }}
                />
              </Box>
              <Text
                pt={2}
                color={'#3C4257'}
                fontSize={'0.8em'}
                fontWeight={isActiveStep(6) ? 'bold' : 'normal'}>
                Submit & Review
              </Text>
            </HStack>
          </Box>
        </Grid>
      </Box>
    </>
  );
};

export default CertificateStepLabel;
