import React, { FC, useEffect } from 'react';
import { Stack, Box, Text, Heading, Table, Tbody, Tr, Td, Button } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { getStepData } from 'utils/utils';

interface BasicDetailsReviewProps {}

const BasicDetailsReview = (props: BasicDetailsReviewProps) => {
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [basicDetail, setBasicDetail] = React.useState<any>({});

  useEffect(() => {
    const stepData = getStepData(steps, 1);
    if (stepData) {
      setBasicDetail(stepData);
    }
  }, [steps]);
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      bg={'white'}
      color={'#252733'}
      maxHeight={367}
      fontSize={18}
      p={5}
      px={5}>
      <Stack width={'100%'}>
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={5}>
          <Heading fontSize={24}>Review 1: Basic Details</Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            height={'34px'}
            _hover={{
              bg: '#10aaed'
            }}>
            {' '}
            Edit{' '}
          </Button>
        </Box>
        <Stack fontSize={18}>
          <Table sx={{ 'td:nth-child(2),td:nth-child(3)': { fontWeight: 'bold' } }}>
            <Tbody>
              <Tr>
                <Td>Website</Td>
                <Td>{basicDetail.website}</Td>
                <Td></Td>
              </Tr>
              <Tr borderStyle={'hidden'}>
                <Td>Date of Incorporation/ Establishmend</Td>
                <Td>{basicDetail.established_on}</Td>
                <Td></Td>
              </Tr>
              <Tr borderStyle={'hidden'}>
                <Td>VASP Category</Td>
                <Td>{basicDetail.vasp_categories}</Td>
                <Td></Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Box>
  );
};
BasicDetailsReview.defaultProps = {
  data: {}
};
export default BasicDetailsReview;
