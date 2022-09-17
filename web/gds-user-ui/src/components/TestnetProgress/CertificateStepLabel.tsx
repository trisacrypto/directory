import { FC } from 'react';
import { Box, Icon, Text, Heading, Stack, Tooltip, Flex } from '@chakra-ui/react';
import { FaCheckCircle, FaDotCircle, FaRegCircle } from 'react-icons/fa';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep } from 'application/store/stepper.slice';
import { findStepKey } from 'utils/utils';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';

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
    <Box
      boxShadow="0px 10px 15px -3px rgba(0,0,0,0.1)"
      borderColor={'#C1C9D2'}
      borderRadius={8}
      borderWidth={1}
      // bg={'white'}
      p={5}
      fontFamily={'Open Sans'}
      width="100%">
      <Heading fontSize={['md', '2xl']} mb={4}>
        <Trans id="Certificate Progress">Certificate Progress</Trans>{' '}
      </Heading>
      <Stack direction="row" spacing={2}>
        <Tooltip
          label={getLabel(1)?.hasError && t`Missing required element`}
          placement="top"
          bg={'red'}>
          <Stack spacing={1} width="100%">
            <Box h="1" borderRadius={'50px'} bg={getLabel(1)?.color} width={'100%'} key={1} />
            <Stack
              direction={{ base: 'column', md: 'row' }}
              alignItems={['center']}
              spacing={{ base: 0, md: 1 }}>
              <Box>
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
                fontWeight={isActiveStep(1) ? 'bold' : 'normal'}
                fontSize={'sm'}
                textAlign="center">
                1 <Trans id="Basic Details">Basic Details</Trans>
              </Text>
            </Stack>
          </Stack>
        </Tooltip>

        <Stack spacing={1} width="100%">
          <Box h="1" bg={getLabel(2)?.color} borderRadius={'50px'} width={'100%'} />
          <Stack
            direction={{ base: 'column', md: 'row' }}
            alignItems={'center'}
            spacing={{ base: 0, md: 1 }}>
            <Box>
              <Icon
                as={getLabel(2)?.icon}
                sx={{
                  path: {
                    fill: getLabel(2)?.color
                  },
                  verticalAlign: 'middle'
                }}
              />
            </Box>
            <Text
              fontSize={'sm'}
              fontWeight={isActiveStep(2) ? 'bold' : 'normal'}
              textAlign="center">
              2 <Trans id="Legal Person">Legal Person</Trans>
            </Text>
          </Stack>
        </Stack>

        <Stack spacing={1} width="100%">
          <Box h="1" bg={getLabel(3)?.color} width={'100%'} borderRadius={'50px'} />
          <Stack
            direction={{ base: 'column', md: 'row' }}
            alignItems={['center']}
            spacing={{ base: 0, md: 1 }}>
            <Box>
              <Icon
                as={getLabel(3)?.icon}
                sx={{
                  path: {
                    fill: getLabel(3)?.color
                  }
                }}
              />
            </Box>
            <Text
              fontSize={'sm'}
              fontWeight={isActiveStep(3) ? 'bold' : 'normal'}
              textAlign="center">
              3 <Trans id="Contacts">Contacts</Trans>
            </Text>
          </Stack>
        </Stack>

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
              fontSize={'sm'}
              fontWeight={isActiveStep(4) ? 'bold' : 'normal'}
              textAlign="center">
              4 <Trans id="TRISA implementation">TRISA implementation</Trans>
            </Text>
          </Stack>
        </Stack>

        <Stack spacing={1} width="100%">
          <Box h="1" bg={getLabel(5)?.color} width={'100%'} borderRadius={'50px'} />
          <Stack
            direction={{ base: 'column', md: 'row' }}
            alignItems={['center']}
            spacing={{ base: 0, md: 1 }}>
            <Box>
              <Icon
                as={getLabel(5)?.icon}
                sx={{
                  path: {
                    fill: getLabel(5)?.color
                  }
                }}
              />
            </Box>
            <Text
              fontSize={'sm'}
              fontWeight={isActiveStep(5) ? 'bold' : 'normal'}
              textAlign="center">
              5 <Trans id="TRIXO Questionnaire">TRIXO Questionnaire</Trans>
            </Text>
          </Stack>
        </Stack>

        <Stack spacing={1} width="100%">
          <Box h="1" bg={getLabel(6)?.color} width={'100%'} borderRadius={'50px'} />
          <Stack
            direction={{ base: 'column', md: 'row' }}
            alignItems={['center']}
            spacing={{ base: 0, md: 1 }}>
            <Box>
              <Icon
                as={getLabel(6)?.icon}
                sx={{
                  path: {
                    fill: getLabel(6)?.color
                  }
                }}
              />
            </Box>
            <Text
              fontSize={'sm'}
              fontWeight={isActiveStep(6) ? 'bold' : 'normal'}
              textAlign="center">
              6 <Trans id="Review">Review</Trans>
            </Text>
          </Stack>
        </Stack>
      </Stack>
      {/* </Grid> */}
    </Box>
  );
};

export default CertificateStepLabel;
