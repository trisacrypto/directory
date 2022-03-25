import React, { FC, useEffect } from 'react';
import { HStack, Box, Icon, Text, Heading, Stack, Grid, Button } from '@chakra-ui/react';
import { FaCheckCircle, FaDotCircle, FaRegCircle } from 'react-icons/fa';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { addStep, setCurrentStep, setStepStatus, TStep } from 'application/store/stepper.slice';
enum LCOLOR {
  'COMPLETE' = '#34A853',
  'PROGRESS' = '#5469D4',
  'SAVE' = '#F29C36',
  'INCOMPLETE' = '#C1C9D2'
}
enum LSTATUS {
  'COMPLETE' = 'complete',
  'PROGRESS' = 'progress',
  'SAVE' = 'save',
  'INCOMPLETE' = 'incomplete'
}
type StepLabelProps = {};

const CertificateStepLabel: FC<StepLabelProps> = (props) => {
  const dispatch = useDispatch();
  const CurrentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const Steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const getStep = (step: number) => {
    return Steps?.filter((s) => s?.key === step);
  };

  // this function need some clean up
  const getLabel = (step: number) => {
    const s = getStep(step);
    if (s && s?.length === 1) {
      if (s[0]?.status === LSTATUS.COMPLETE) {
        return {
          color: LCOLOR.COMPLETE,
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
          <Heading fontSize={20}> Progress bar </Heading>
        </Box>
        <Grid templateColumns="repeat(6, 1fr)" gap={2}>
          <Box w="70px" h="1" borderRadius={50} bg={getLabel(1)?.color} width={'100%'} key={1}>
            <HStack>
              <Box pt={3}>
                <Icon as={getLabel(1)?.icon} color={getLabel(1)?.color} />
              </Box>
              <Text pt={2} color={'#3C4257'} fontSize={'0.8em'}>
                Basic Details
              </Text>
            </HStack>
          </Box>

          <Box w="70px" h="1" bg={getLabel(2)?.color} width={'100%'}>
            <HStack>
              <Box pt={3}>
                <Icon as={getLabel(2)?.icon} color={getLabel(2)?.color} />
              </Box>
              <Text pt={2} color={'#3C4257'} fontSize={'0.8em'}>
                Legal Person
              </Text>
            </HStack>
          </Box>

          <Box w="70px" h="1" bg={getLabel(3)?.color} width={'100%'}>
            <HStack>
              <Box pt={3}>
                <Icon as={getLabel(3)?.icon} color={getLabel(3)?.color} />
              </Box>
              <Text pt={2} color={'#3C4257'} fontSize={'0.8em'}>
                Contacts
              </Text>
            </HStack>
          </Box>

          <Box w="70px" h="1" bg={getLabel(4)?.color} width={'100%'}>
            <HStack>
              <Box pt={3}>
                <Icon as={getLabel(4)?.icon} color={getLabel(4)?.color} />
              </Box>
              <Text pt={2} color={'#3C4257'} fontSize={'0.8em'}>
                Trisa implementation
              </Text>
            </HStack>
          </Box>

          <Box w="70px" h="1" bg={getLabel(5)?.color} width={'100%'}>
            <HStack>
              <Box pt={3}>
                <Icon as={getLabel(5)?.icon} color={getLabel(5)?.color} />
              </Box>
              <Text pt={2} color={'#3C4257'} fontSize={'0.8em'}>
                TRIXO Questionnaire
              </Text>
            </HStack>
          </Box>

          <Box w="70px" h="1" bg={getLabel(6)?.color} width={'100%'}>
            <HStack>
              <Box pt={3}>
                <Icon as={getLabel(6)?.icon} color={getLabel(6)?.color} />
              </Box>
              <Text pt={2} color={'#3C4257'} fontSize={'0.8em'}>
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
